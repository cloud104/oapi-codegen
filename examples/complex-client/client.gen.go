package examples

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
	yaml "go.yaml.in/yaml/v3"
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
	listBeers(ctx context.Context, params *ListBeersParams, reqEditors ...RequestEditorFn) (*http.Response, error)
	createBeer(ctx context.Context, body CreateBeerJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
	createBeerWithFormdataBody(ctx context.Context, body CreateBeerFormdataRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
	createBeerWithTextBody(ctx context.Context, body CreateBeerTextRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
	exportBeers(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
	deleteBeer(ctx context.Context, beerID int, reqEditors ...RequestEditorFn) (*http.Response, error)
	getBeer(ctx context.Context, beerID int, reqEditors ...RequestEditorFn) (*http.Response, error)
	updateBeer(ctx context.Context, beerID int, body UpdateBeerJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
	replaceBeer(ctx context.Context, beerID int, body ReplaceBeerJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *rawClient) listBeers(ctx context.Context, params *ListBeersParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newListBeersRequest(c.Server, params)
	if err != nil {
		return nil, fmt.Errorf("create GET /beers request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /beers request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /beers request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) createBeer(ctx context.Context, body CreateBeerJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newCreateBeerRequest(c.Server, body)
	if err != nil {
		return nil, fmt.Errorf("create POST /beers request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to POST /beers request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform POST /beers request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) createBeerWithFormdataBody(ctx context.Context, body CreateBeerFormdataRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newCreateBeerRequestWithFormdataBody(c.Server, body)
	if err != nil {
		return nil, fmt.Errorf("create POST /beers request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to POST /beers request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform POST /beers request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) createBeerWithTextBody(ctx context.Context, body CreateBeerTextRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newCreateBeerRequestWithTextBody(c.Server, body)
	if err != nil {
		return nil, fmt.Errorf("create POST /beers request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to POST /beers request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform POST /beers request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) exportBeers(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newExportBeersRequest(c.Server)
	if err != nil {
		return nil, fmt.Errorf("create GET /beers/export request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /beers/export request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /beers/export request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) deleteBeer(ctx context.Context, beerID int, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newDeleteBeerRequest(c.Server, beerID)
	if err != nil {
		return nil, fmt.Errorf("create DELETE /beers/{beerId} request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to DELETE /beers/{beerId} request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform DELETE /beers/{beerId} request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) getBeer(ctx context.Context, beerID int, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newGetBeerRequest(c.Server, beerID)
	if err != nil {
		return nil, fmt.Errorf("create GET /beers/{beerId} request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /beers/{beerId} request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /beers/{beerId} request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) updateBeer(ctx context.Context, beerID int, body UpdateBeerJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newUpdateBeerRequest(c.Server, beerID, body)
	if err != nil {
		return nil, fmt.Errorf("create PATCH /beers/{beerId} request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to PATCH /beers/{beerId} request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform PATCH /beers/{beerId} request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) replaceBeer(ctx context.Context, beerID int, body ReplaceBeerJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newReplaceBeerRequest(c.Server, beerID, body)
	if err != nil {
		return nil, fmt.Errorf("create PUT /beers/{beerId} request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to PUT /beers/{beerId} request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform PUT /beers/{beerId} request: %w", err)
	}

	return resp, nil
}

func newListBeersRequest(server string, params *ListBeersParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/beers"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()
		var rawQueryFragments []string

		if params.Page != nil {
			if queryFrag, err := runtime.StyleParamWithOptions("form", true, "page", *params.Page, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationQuery, Type: "integer", Format: ""}); err != nil {
				return nil, err
			} else {
				for _, qp := range strings.Split(queryFrag, "&") {
					rawQueryFragments = append(rawQueryFragments, qp)
				}
			}
		}

		if params.PageSize != nil {
			if queryFrag, err := runtime.StyleParamWithOptions("form", true, "pageSize", *params.PageSize, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationQuery, Type: "integer", Format: ""}); err != nil {
				return nil, err
			} else {
				for _, qp := range strings.Split(queryFrag, "&") {
					rawQueryFragments = append(rawQueryFragments, qp)
				}
			}
		}

		if encoded := queryValues.Encode(); encoded != "" {
			rawQueryFragments = append(rawQueryFragments, encoded)
		}

		queryURL.RawQuery = strings.Join(rawQueryFragments, "&")
	}

	req, err := http.NewRequest(http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func newCreateBeerRequest(server string, body CreateBeerJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader

	buf, err := runtime.JSONMarshal(body)
	if err != nil {
		return nil, err
	}

	bodyReader = bytes.NewReader(buf)

	return newCreateBeerRequestWithBody(server, "application/json", bodyReader)
}

func newCreateBeerRequestWithFormdataBody(server string, body CreateBeerFormdataRequestBody) (*http.Request, error) {
	var bodyReader io.Reader

	bodyStr, err := runtime.MarshalForm(body, nil)
	if err != nil {
		return nil, err
	}

	bodyReader = strings.NewReader(bodyStr.Encode())

	return newCreateBeerRequestWithBody(server, "application/x-www-form-urlencoded", bodyReader)
}

func newCreateBeerRequestWithTextBody(server string, body CreateBeerTextRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	if stringer, ok := any(body).(fmt.Stringer); ok {
		bodyReader = strings.NewReader(stringer.String())
	} else {
		bodyReader = strings.NewReader(fmt.Sprint(body))
	}

	return newCreateBeerRequestWithBody(server, "text/plain", bodyReader)
}

func newCreateBeerRequestWithBody(server, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/beers"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func newExportBeersRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/beers/export"
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

func newDeleteBeerRequest(server string, beerID int) (*http.Request, error) {
	var err error
	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "beerId", beerID, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationPath, Type: "integer", Format: ""})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/beers/" + pathParam0
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodDelete, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func newGetBeerRequest(server string, beerID int) (*http.Request, error) {
	var err error
	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "beerId", beerID, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationPath, Type: "integer", Format: ""})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/beers/" + pathParam0
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

func newUpdateBeerRequest(server string, beerID int, body UpdateBeerJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader

	buf, err := runtime.JSONMarshal(body)
	if err != nil {
		return nil, err
	}

	bodyReader = bytes.NewReader(buf)

	return newUpdateBeerRequestWithBody(server, beerID, "application/json", bodyReader)
}

func newUpdateBeerRequestWithBody(server string, beerID int, contentType string, body io.Reader) (*http.Request, error) {
	var err error
	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "beerId", beerID, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationPath, Type: "integer", Format: ""})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/beers/" + pathParam0
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func newReplaceBeerRequest(server string, beerID int, body ReplaceBeerJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader

	buf, err := runtime.JSONMarshal(body)
	if err != nil {
		return nil, err
	}

	bodyReader = bytes.NewReader(buf)

	return newReplaceBeerRequestWithBody(server, beerID, "application/json", bodyReader)
}

func newReplaceBeerRequestWithBody(server string, beerID int, contentType string, body io.Reader) (*http.Request, error) {
	var err error
	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "beerId", beerID, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationPath, Type: "integer", Format: ""})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/beers/" + pathParam0
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

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
	ListBeers(ctx context.Context, params *ListBeersParams, reqEditors ...RequestEditorFn) (*ListBeersResponse, error)
	CreateBeer(ctx context.Context, body CreateBeerJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateBeerResponse, error)
	CreateBeerWithFormdataBody(ctx context.Context, body CreateBeerFormdataRequestBody, reqEditors ...RequestEditorFn) (*CreateBeerResponse, error)
	CreateBeerWithTextBody(ctx context.Context, body CreateBeerTextRequestBody, reqEditors ...RequestEditorFn) (*CreateBeerResponse, error)
	ExportBeers(ctx context.Context, reqEditors ...RequestEditorFn) (*ExportBeersResponse, error)
	DeleteBeer(ctx context.Context, beerID int, reqEditors ...RequestEditorFn) (*DeleteBeerResponse, error)
	GetBeer(ctx context.Context, beerID int, reqEditors ...RequestEditorFn) (*GetBeerResponse, error)
	UpdateBeer(ctx context.Context, beerID int, body UpdateBeerJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateBeerResponse, error)
	ReplaceBeer(ctx context.Context, beerID int, body ReplaceBeerJSONRequestBody, reqEditors ...RequestEditorFn) (*ReplaceBeerResponse, error)
}

type ListBeersResponse200Headers struct {
	Link        *string
	XTotalCount *int
}

type ListBeersResponse202Headers struct {
	Link        *string
	RetryAfter  *int
	XRequestID  *string
	XTotalCount int
}

type ListBeersResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *BeerPage
	XML200       *BeerPage
	XML202       *BeerPage
	Headers200   *ListBeersResponse200Headers
	Headers202   *ListBeersResponse202Headers
}

func (r ListBeersResponse) GetBeerPage() (*BeerPage, error) {
	if r.JSON200 != nil {
		return r.JSON200, nil
	}

	if r.XML200 != nil {
		return r.XML200, nil
	}

	if r.XML202 != nil {
		return r.XML202, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 200, 202:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r ListBeersResponse) GetBody() []byte {
	return r.Body
}

func (r ListBeersResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r ListBeersResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r ListBeersResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

func (r ListBeersResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

type CreateBeerResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Beer
	JSON400      *BadRequest
	JSON409      *Conflict
	JSON415      *UnsupportedMediaType
}

func (r CreateBeerResponse) GetBeer() (*Beer, error) {
	if r.JSON201 != nil {
		return r.JSON201, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 201:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r CreateBeerResponse) GetBadRequest() (*BadRequest, error) {
	if r.JSON400 != nil {
		return r.JSON400, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 400:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r CreateBeerResponse) GetConflict() (*Conflict, error) {
	if r.JSON409 != nil {
		return r.JSON409, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 409:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r CreateBeerResponse) GetUnsupportedMediaType() (*UnsupportedMediaType, error) {
	if r.JSON415 != nil {
		return r.JSON415, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 415:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r CreateBeerResponse) GetBody() []byte {
	return r.Body
}

func (r CreateBeerResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r CreateBeerResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r CreateBeerResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

func (r CreateBeerResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

type ExportBeersResponse200Headers struct {
	ContentDisposition string
}

type ExportBeersResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	YAML200      *openapi_types.File
	JSON400      *BadRequest
	Headers200   *ExportBeersResponse200Headers
}

func (r ExportBeersResponse) GetFile() (*openapi_types.File, error) {
	if r.YAML200 != nil {
		return r.YAML200, nil
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

func (r ExportBeersResponse) GetBadRequest() (*BadRequest, error) {
	if r.JSON400 != nil {
		return r.JSON400, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 400:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r ExportBeersResponse) GetBody() []byte {
	return r.Body
}

func (r ExportBeersResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r ExportBeersResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r ExportBeersResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

func (r ExportBeersResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

type DeleteBeerResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON404      *NotFound
	JSON409      *Conflict
	JSONDefault  *map[string]any
}

func (r DeleteBeerResponse) GetNotFound() (*NotFound, error) {
	if r.JSON404 != nil {
		return r.JSON404, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 404:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r DeleteBeerResponse) GetConflict() (*Conflict, error) {
	if r.JSON409 != nil {
		return r.JSON409, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 409:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r DeleteBeerResponse) GetBody() []byte {
	return r.Body
}

func (r DeleteBeerResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r DeleteBeerResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r DeleteBeerResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

func (r DeleteBeerResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

type GetBeerResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Abv   *float32 `json:"abv,omitempty"`
		ID    int      `json:"id"`
		Name  string   `json:"name"`
		Style *string  `json:"style,omitempty"`
	}
	JSON404 *NotFound
}

func (r GetBeerResponse) GetNotFound() (*NotFound, error) {
	if r.JSON404 != nil {
		return r.JSON404, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 404:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r GetBeerResponse) GetBody() []byte {
	return r.Body
}

func (r GetBeerResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r GetBeerResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r GetBeerResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

func (r GetBeerResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

type UpdateBeerResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Beer
	JSON400      *BadRequest
	JSON404      *NotFound
	JSON409      *Conflict
	JSON415      *UnsupportedMediaType
}

func (r UpdateBeerResponse) GetBeer() (*Beer, error) {
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

func (r UpdateBeerResponse) GetBadRequest() (*BadRequest, error) {
	if r.JSON400 != nil {
		return r.JSON400, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 400:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r UpdateBeerResponse) GetNotFound() (*NotFound, error) {
	if r.JSON404 != nil {
		return r.JSON404, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 404:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r UpdateBeerResponse) GetConflict() (*Conflict, error) {
	if r.JSON409 != nil {
		return r.JSON409, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 409:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r UpdateBeerResponse) GetUnsupportedMediaType() (*UnsupportedMediaType, error) {
	if r.JSON415 != nil {
		return r.JSON415, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 415:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r UpdateBeerResponse) GetBody() []byte {
	return r.Body
}

func (r UpdateBeerResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r UpdateBeerResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r UpdateBeerResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

func (r UpdateBeerResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

type ReplaceBeerResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Beer
	JSON400      *BadRequest
	JSON404      *NotFound
	JSON409      *Conflict
	JSON415      *UnsupportedMediaType
}

func (r ReplaceBeerResponse) GetBeer() (*Beer, error) {
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

func (r ReplaceBeerResponse) GetBadRequest() (*BadRequest, error) {
	if r.JSON400 != nil {
		return r.JSON400, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 400:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r ReplaceBeerResponse) GetNotFound() (*NotFound, error) {
	if r.JSON404 != nil {
		return r.JSON404, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 404:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r ReplaceBeerResponse) GetConflict() (*Conflict, error) {
	if r.JSON409 != nil {
		return r.JSON409, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 409:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r ReplaceBeerResponse) GetUnsupportedMediaType() (*UnsupportedMediaType, error) {
	if r.JSON415 != nil {
		return r.JSON415, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 415:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r ReplaceBeerResponse) GetBody() []byte {
	return r.Body
}

func (r ReplaceBeerResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r ReplaceBeerResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r ReplaceBeerResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

func (r ReplaceBeerResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

func (c *Client) ListBeers(ctx context.Context, params *ListBeersParams, reqEditors ...RequestEditorFn) (*ListBeersResponse, error) {
	rsp, err := c.listBeers(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseListBeersResponse(rsp)
}

func (c *Client) CreateBeer(ctx context.Context, body CreateBeerJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateBeerResponse, error) {
	rsp, err := c.createBeer(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseCreateBeerResponse(rsp)
}

func (c *Client) CreateBeerWithFormdataBody(ctx context.Context, body CreateBeerFormdataRequestBody, reqEditors ...RequestEditorFn) (*CreateBeerResponse, error) {
	rsp, err := c.createBeerWithFormdataBody(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseCreateBeerResponse(rsp)
}

func (c *Client) CreateBeerWithTextBody(ctx context.Context, body CreateBeerTextRequestBody, reqEditors ...RequestEditorFn) (*CreateBeerResponse, error) {
	rsp, err := c.createBeerWithTextBody(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseCreateBeerResponse(rsp)
}

func (c *Client) ExportBeers(ctx context.Context, reqEditors ...RequestEditorFn) (*ExportBeersResponse, error) {
	rsp, err := c.exportBeers(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseExportBeersResponse(rsp)
}

func (c *Client) DeleteBeer(ctx context.Context, beerID int, reqEditors ...RequestEditorFn) (*DeleteBeerResponse, error) {
	rsp, err := c.deleteBeer(ctx, beerID, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseDeleteBeerResponse(rsp)
}

func (c *Client) GetBeer(ctx context.Context, beerID int, reqEditors ...RequestEditorFn) (*GetBeerResponse, error) {
	rsp, err := c.getBeer(ctx, beerID, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseGetBeerResponse(rsp)
}

func (c *Client) UpdateBeer(ctx context.Context, beerID int, body UpdateBeerJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateBeerResponse, error) {
	rsp, err := c.updateBeer(ctx, beerID, body, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseUpdateBeerResponse(rsp)
}

func (c *Client) ReplaceBeer(ctx context.Context, beerID int, body ReplaceBeerJSONRequestBody, reqEditors ...RequestEditorFn) (*ReplaceBeerResponse, error) {
	rsp, err := c.replaceBeer(ctx, beerID, body, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseReplaceBeerResponse(rsp)
}

func parseListBeersResponse(rsp *http.Response) (*ListBeersResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &ListBeersResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}
	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest BeerPage
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON200 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "xml") && rsp.StatusCode == 200:
		var dest BeerPage
		if err := xml.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.XML200 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "xml") && rsp.StatusCode == 202:
		var dest BeerPage
		if err := xml.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.XML202 = &dest
	}

	switch {
	case rsp.StatusCode == 200:
		var headers ListBeersResponse200Headers

		if values := rsp.Header.Values("Link"); len(values) > 0 {
			var value string
			if err := runtime.BindStyledParameterWithOptions("simple", "Link", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "string", Format: ""}); err != nil {
				return nil, err
			}

			headers.Link = &value
		}

		if values := rsp.Header.Values("X-Total-Count"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "X-Total-Count", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}

			headers.XTotalCount = &value
		}

		response.Headers200 = &headers
	case rsp.StatusCode == 202:
		var headers ListBeersResponse202Headers

		if values := rsp.Header.Values("Link"); len(values) > 0 {
			var value string
			if err := runtime.BindStyledParameterWithOptions("simple", "Link", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "string", Format: ""}); err != nil {
				return nil, err
			}

			headers.Link = &value
		}

		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}

			headers.RetryAfter = &value
		}

		if values := rsp.Header.Values("X-Request-Id"); len(values) > 0 {
			var value string
			if err := runtime.BindStyledParameterWithOptions("simple", "X-Request-Id", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "string", Format: ""}); err != nil {
				return nil, err
			}

			headers.XRequestID = &value
		}

		if values := rsp.Header.Values("X-Total-Count"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "X-Total-Count", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}

			headers.XTotalCount = value
		}

		response.Headers202 = &headers
	}

	return response, nil
}

func parseCreateBeerResponse(rsp *http.Response) (*CreateBeerResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &CreateBeerResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}
	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest Beer
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON201 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON400 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 409:
		var dest Conflict
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON409 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 415:
		var dest UnsupportedMediaType
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON415 = &dest
	}

	return response, nil
}

func parseExportBeersResponse(rsp *http.Response) (*ExportBeersResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &ExportBeersResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}
	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "yaml") && rsp.StatusCode == 200:
		var dest openapi_types.File
		if err := yaml.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.YAML200 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON400 = &dest
	}

	switch {
	case rsp.StatusCode == 200:
		var headers ExportBeersResponse200Headers

		if values := rsp.Header.Values("Content-Disposition"); len(values) > 0 {
			var value string
			if err := runtime.BindStyledParameterWithOptions("simple", "Content-Disposition", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true, Type: "string", Format: ""}); err != nil {
				return nil, err
			}

			headers.ContentDisposition = value
		}

		response.Headers200 = &headers
	}

	return response, nil
}

func parseDeleteBeerResponse(rsp *http.Response) (*DeleteBeerResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &DeleteBeerResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}
	switch {
	case rsp.StatusCode == 204:
		break
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON404 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 409:
		var dest Conflict
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON409 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest map[string]any
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSONDefault = &dest
	}

	return response, nil
}

func parseGetBeerResponse(rsp *http.Response) (*GetBeerResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &GetBeerResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}
	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Abv   *float32 `json:"abv,omitempty"`
			ID    int      `json:"id"`
			Name  string   `json:"name"`
			Style *string  `json:"style,omitempty"`
		}
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON200 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON404 = &dest
	}

	return response, nil
}

func parseUpdateBeerResponse(rsp *http.Response) (*UpdateBeerResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &UpdateBeerResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}
	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Beer
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON200 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON400 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON404 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 409:
		var dest Conflict
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON409 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 415:
		var dest UnsupportedMediaType
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON415 = &dest
	}

	return response, nil
}

func parseReplaceBeerResponse(rsp *http.Response) (*ReplaceBeerResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &ReplaceBeerResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}
	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Beer
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON200 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON400 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON404 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 409:
		var dest Conflict
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON409 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 415:
		var dest UnsupportedMediaType
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON415 = &dest
	}

	return response, nil
}
