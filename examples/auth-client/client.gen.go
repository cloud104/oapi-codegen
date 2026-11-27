package examples
import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"github.com/oapi-codegen/runtime"
)
type RequestEditorFn func(ctx context.Context, req *http.Request) error
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
type rawClient struct {
	Server string
	Client HttpRequestDoer
	RequestEditors []RequestEditorFn
	BasicAuth struct {
		Username string
		Password string
	}
	BearerAuth struct {
		Provider BearerTokenProvider
	}
	APIKeyAuth struct {
		APIKey string
	}
}
type BearerTokenProvider interface {
	Token(ctx context.Context) (string, error)
	Invalidate(ctx context.Context) error
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
func WithBasicAuth(username, password string) ClientOption {
	return func(c *rawClient) error {
		c.BasicAuth.Username = username
		c.BasicAuth.Password = password
		return nil
	}
}
func WithBearerAuth(provider BearerTokenProvider) ClientOption {
	return func(c *rawClient) error {
		c.BearerAuth.Provider = provider
		return nil
	}
}
func WithAPIKeyAuth(apiKey string) ClientOption {
	return func(c *rawClient) error {
		c.APIKeyAuth.APIKey = apiKey
		return nil
	}
}
type rawClientInterface interface {
	heartbeat(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
	ping(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}
func (c *rawClient) heartbeat(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newHeartbeatRequest(c.Server)
	if err != nil {
		return nil, fmt.Errorf("create GET /heartbeat request: %w", err)
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /heartbeat request: %w", err)
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /heartbeat request: %w", err)
	}
	return resp, nil
}
func (c *rawClient) ping(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	editors := append(slices.Clone(reqEditors), func(ctx context.Context, req *http.Request) error {
		switch {
		case c.BasicAuth.Username != "" && c.BasicAuth.Password != "":
			req.SetBasicAuth(c.BasicAuth.Username, c.BasicAuth.Password)
		case c.BearerAuth.Provider != nil:
			token, err := c.BearerAuth.Provider.Token(ctx)
			if err != nil {
				return fmt.Errorf("get authentication token for GET /ping: %w", err)
			}
			req.Header.Set("Authorization", "Bearer "+token)
		case c.APIKeyAuth.APIKey != "":
			q := req.URL.Query()
			q.Set("X-API-Key", c.APIKeyAuth.APIKey)
			req.URL.RawQuery = q.Encode()
		}
		return nil
	})
	rsp, err := c.unauthenticatedPing(ctx, editors...)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != http.StatusUnauthorized {
		return rsp, nil
	}
	_ = rsp.Body.Close()
	if provider := c.BearerAuth.Provider; provider != nil {
		if err := provider.Invalidate(ctx); err != nil {
			return nil, fmt.Errorf("invalidate token after unauthorized GET /ping response: %w", err)
		}
	}
	return c.unauthenticatedPing(ctx, editors...)
}
func (c *rawClient) unauthenticatedPing(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newPingRequest(c.Server)
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
func newHeartbeatRequest(server string) (*http.Request, error) {
	var err error
	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}
	operationPath := "/heartbeat"
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
func newPingRequest(server string) (*http.Request, error) {
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
	Heartbeat(ctx context.Context, reqEditors ...RequestEditorFn) (*HeartbeatResponse, error)
	Ping(ctx context.Context, reqEditors ...RequestEditorFn) (*PingResponse, error)
}
type HeartbeatResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200 *Pong
}
func (r HeartbeatResponse) GetPong() (*Pong, error) {
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
func (r HeartbeatResponse) GetBody() []byte {
	return r.Body
}
func (r HeartbeatResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}
	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}
func (r HeartbeatResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}
func (r HeartbeatResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}
func (r HeartbeatResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}
type PingResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200 *Pong
}
func (r PingResponse) GetPong() (*Pong, error) {
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
func (r PingResponse) GetBody() []byte {
	return r.Body
}
func (r PingResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}
	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}
func (r PingResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}
func (r PingResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}
func (r PingResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}
func (c *Client) Heartbeat(ctx context.Context, reqEditors ...RequestEditorFn) (*HeartbeatResponse, error) {
	rsp, err := c.heartbeat(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return parseHeartbeatResponse(rsp)
}
func (c *Client) Ping(ctx context.Context, reqEditors ...RequestEditorFn) (*PingResponse, error) {
	rsp, err := c.ping(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return parsePingResponse(rsp)
}
func parseHeartbeatResponse(rsp *http.Response) (*HeartbeatResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}
	response := &HeartbeatResponse{
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
func parsePingResponse(rsp *http.Response) (*PingResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}
	response := &PingResponse{
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
