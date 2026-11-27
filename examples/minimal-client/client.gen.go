package examples

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oapi-codegen/runtime"
)

type RequestEditorFn func(ctx context.Context, req *http.Request) error

type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type rawClient struct {
	Server         string
	Client         HttpRequestDoer
	RequestEditors []RequestEditorFn
}

type ClientOption func(*rawClient) error

func newRawClient(server string, opts ...ClientOption) (*rawClient, error) {
	client := rawClient{
		Server: server,
	}
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}

	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}

	if client.Client == nil {
		client.Client = &http.Client{}
	}

	return &client, nil
}

func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *rawClient) error {
		c.Client = doer
		return nil
	}
}

func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *rawClient) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

type rawClientInterface interface {
	getPing(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *rawClient) getPing(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newGetPingRequest(c.Server)
	if err != nil {
		return nil, fmt.Errorf("create GET /ping request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /ping request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /ping request: %w", err)
	}

	return resp, nil
}

func newGetPingRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/ping"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *rawClient) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}

	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}

	return nil
}

type Client struct {
	rawClientInterface
}

func NewClient(server string, opts ...ClientOption) (*Client, error) {
	client, err := newRawClient(server, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{client}, nil
}

func WithBaseURL(baseURL string) ClientOption {
	return func(c *rawClient) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}

		c.Server = newBaseURL.String()

		return nil
	}
}

type ClientInterface interface {
	GetPing(ctx context.Context, reqEditors ...RequestEditorFn) (*GetPingResponse, error)
}

type GetPingResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Pong
}

func (r GetPingResponse) GetPong() (*Pong, error) {
	if r.JSON200 != nil {
		return r.JSON200, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 200:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r GetPingResponse) GetBody() []byte {
	return r.Body
}

func (r GetPingResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r GetPingResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r GetPingResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

func (r GetPingResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

func (c *Client) GetPing(ctx context.Context, reqEditors ...RequestEditorFn) (*GetPingResponse, error) {
	rsp, err := c.getPing(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseGetPingResponse(rsp)
}

func parseGetPingResponse(rsp *http.Response) (*GetPingResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &GetPingResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}
	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Pong
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON200 = &dest
	}

	return response, nil
}
