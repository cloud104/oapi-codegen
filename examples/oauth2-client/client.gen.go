package examples

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

type RequestEditorFn func(ctx context.Context, req *http.Request) error

type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type rawClient struct {
	Server         string
	Client         HttpRequestDoer
	RequestEditors []RequestEditorFn
	OAuth2Password struct {
		TokenSource oauth2.TokenSource
	}
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

func WithOAuth2Password(clientID, clientSecret, username, password string) ClientOption {
	return func(c *rawClient) error {
		config := &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint: oauth2.Endpoint{
				TokenURL: "https://auth.example.com/oauth2/token",
			},
		}

		token, err := config.PasswordCredentialsToken(context.TODO(), username, password)
		if err != nil {
			return fmt.Errorf("authenticate: %w", err)
		}

		c.OAuth2Password.TokenSource = config.TokenSource(context.TODO(), token)

		return nil
	}
}

type rawClientInterface interface {
	getDevices(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
	getLegacyImplicit(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
	getLegacyPassword(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
	getOptionalAuth(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
	getPublic(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
	getReports(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
	getServices(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
	listUsers(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
	createUser(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *rawClient) getDevices(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newGetDevicesRequest(c.Server)
	if err != nil {
		return nil, fmt.Errorf("create GET /devices request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /devices request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /devices request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) getLegacyImplicit(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newGetLegacyImplicitRequest(c.Server)
	if err != nil {
		return nil, fmt.Errorf("create GET /legacy/implicit request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /legacy/implicit request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /legacy/implicit request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) getLegacyPassword(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newGetLegacyPasswordRequest(c.Server)
	if err != nil {
		return nil, fmt.Errorf("create GET /legacy/password request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /legacy/password request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /legacy/password request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) getOptionalAuth(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newGetOptionalAuthRequest(c.Server)
	if err != nil {
		return nil, fmt.Errorf("create GET /optional-auth request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /optional-auth request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /optional-auth request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) getPublic(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newGetPublicRequest(c.Server)
	if err != nil {
		return nil, fmt.Errorf("create GET /public request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /public request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /public request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) getReports(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newGetReportsRequest(c.Server)
	if err != nil {
		return nil, fmt.Errorf("create GET /reports request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /reports request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /reports request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) getServices(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newGetServicesRequest(c.Server)
	if err != nil {
		return nil, fmt.Errorf("create GET /services request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /services request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /services request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) listUsers(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newListUsersRequest(c.Server)
	if err != nil {
		return nil, fmt.Errorf("create GET /users request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /users request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /users request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) createUser(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newCreateUserRequest(c.Server)
	if err != nil {
		return nil, fmt.Errorf("create POST /users request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to POST /users request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform POST /users request: %w", err)
	}

	return resp, nil
}

func newGetDevicesRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/devices"
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

func newGetLegacyImplicitRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/legacy/implicit"
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

func newGetLegacyPasswordRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/legacy/password"
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

func newGetOptionalAuthRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/optional-auth"
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

func newGetPublicRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/public"
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

func newGetReportsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/reports"
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

func newGetServicesRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/services"
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

func newListUsersRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/users"
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

func newCreateUserRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/users"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, queryURL.String(), nil)
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
	GetDevices(ctx context.Context, reqEditors ...RequestEditorFn) (*GetDevicesResponse, error)
	GetLegacyImplicit(ctx context.Context, reqEditors ...RequestEditorFn) (*GetLegacyImplicitResponse, error)
	GetLegacyPassword(ctx context.Context, reqEditors ...RequestEditorFn) (*GetLegacyPasswordResponse, error)
	GetOptionalAuth(ctx context.Context, reqEditors ...RequestEditorFn) (*GetOptionalAuthResponse, error)
	GetPublic(ctx context.Context, reqEditors ...RequestEditorFn) (*GetPublicResponse, error)
	GetReports(ctx context.Context, reqEditors ...RequestEditorFn) (*GetReportsResponse, error)
	GetServices(ctx context.Context, reqEditors ...RequestEditorFn) (*GetServicesResponse, error)
	ListUsers(ctx context.Context, reqEditors ...RequestEditorFn) (*ListUsersResponse, error)
	CreateUser(ctx context.Context, reqEditors ...RequestEditorFn) (*CreateUserResponse, error)
}

type GetDevicesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

func (r GetDevicesResponse) GetBody() []byte {
	return r.Body
}

func (r GetDevicesResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r GetDevicesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r GetDevicesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

func (r GetDevicesResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

type GetLegacyImplicitResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

func (r GetLegacyImplicitResponse) GetBody() []byte {
	return r.Body
}

func (r GetLegacyImplicitResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r GetLegacyImplicitResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r GetLegacyImplicitResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

func (r GetLegacyImplicitResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

type GetLegacyPasswordResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

func (r GetLegacyPasswordResponse) GetBody() []byte {
	return r.Body
}

func (r GetLegacyPasswordResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r GetLegacyPasswordResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r GetLegacyPasswordResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

func (r GetLegacyPasswordResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

type GetOptionalAuthResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

func (r GetOptionalAuthResponse) GetBody() []byte {
	return r.Body
}

func (r GetOptionalAuthResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r GetOptionalAuthResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r GetOptionalAuthResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

func (r GetOptionalAuthResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

type GetPublicResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

func (r GetPublicResponse) GetBody() []byte {
	return r.Body
}

func (r GetPublicResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r GetPublicResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r GetPublicResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

func (r GetPublicResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

type GetReportsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

func (r GetReportsResponse) GetBody() []byte {
	return r.Body
}

func (r GetReportsResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r GetReportsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r GetReportsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

func (r GetReportsResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

type GetServicesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

func (r GetServicesResponse) GetBody() []byte {
	return r.Body
}

func (r GetServicesResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r GetServicesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r GetServicesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

func (r GetServicesResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

type ListUsersResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

func (r ListUsersResponse) GetBody() []byte {
	return r.Body
}

func (r ListUsersResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r ListUsersResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r ListUsersResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

func (r ListUsersResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

type CreateUserResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

func (r CreateUserResponse) GetBody() []byte {
	return r.Body
}

func (r CreateUserResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r CreateUserResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r CreateUserResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

func (r CreateUserResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

func (c *Client) GetDevices(ctx context.Context, reqEditors ...RequestEditorFn) (*GetDevicesResponse, error) {
	rsp, err := c.getDevices(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseGetDevicesResponse(rsp)
}

func (c *Client) GetLegacyImplicit(ctx context.Context, reqEditors ...RequestEditorFn) (*GetLegacyImplicitResponse, error) {
	rsp, err := c.getLegacyImplicit(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseGetLegacyImplicitResponse(rsp)
}

func (c *Client) GetLegacyPassword(ctx context.Context, reqEditors ...RequestEditorFn) (*GetLegacyPasswordResponse, error) {
	rsp, err := c.getLegacyPassword(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseGetLegacyPasswordResponse(rsp)
}

func (c *Client) GetOptionalAuth(ctx context.Context, reqEditors ...RequestEditorFn) (*GetOptionalAuthResponse, error) {
	rsp, err := c.getOptionalAuth(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseGetOptionalAuthResponse(rsp)
}

func (c *Client) GetPublic(ctx context.Context, reqEditors ...RequestEditorFn) (*GetPublicResponse, error) {
	rsp, err := c.getPublic(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseGetPublicResponse(rsp)
}

func (c *Client) GetReports(ctx context.Context, reqEditors ...RequestEditorFn) (*GetReportsResponse, error) {
	rsp, err := c.getReports(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseGetReportsResponse(rsp)
}

func (c *Client) GetServices(ctx context.Context, reqEditors ...RequestEditorFn) (*GetServicesResponse, error) {
	rsp, err := c.getServices(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseGetServicesResponse(rsp)
}

func (c *Client) ListUsers(ctx context.Context, reqEditors ...RequestEditorFn) (*ListUsersResponse, error) {
	rsp, err := c.listUsers(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseListUsersResponse(rsp)
}

func (c *Client) CreateUser(ctx context.Context, reqEditors ...RequestEditorFn) (*CreateUserResponse, error) {
	rsp, err := c.createUser(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseCreateUserResponse(rsp)
}

func parseGetDevicesResponse(rsp *http.Response) (*GetDevicesResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &GetDevicesResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

func parseGetLegacyImplicitResponse(rsp *http.Response) (*GetLegacyImplicitResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &GetLegacyImplicitResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

func parseGetLegacyPasswordResponse(rsp *http.Response) (*GetLegacyPasswordResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &GetLegacyPasswordResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

func parseGetOptionalAuthResponse(rsp *http.Response) (*GetOptionalAuthResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &GetOptionalAuthResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

func parseGetPublicResponse(rsp *http.Response) (*GetPublicResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &GetPublicResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

func parseGetReportsResponse(rsp *http.Response) (*GetReportsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &GetReportsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

func parseGetServicesResponse(rsp *http.Response) (*GetServicesResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &GetServicesResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

func parseListUsersResponse(rsp *http.Response) (*ListUsersResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &ListUsersResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

func parseCreateUserResponse(rsp *http.Response) (*CreateUserResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &CreateUserResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}
