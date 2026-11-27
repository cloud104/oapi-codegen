package examples

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/oapi-codegen/runtime"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func newGetHealthRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/health"
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

type ClientInterface interface {
	StartChargingSession(ctx context.Context, params *StartChargingSessionParams, body StartChargingSessionJSONRequestBody, reqEditors ...RequestEditorFn) (*StartChargingSessionResponse, error)
	GetChargingSession(ctx context.Context, sessionID SessionID, reqEditors ...RequestEditorFn) (*GetChargingSessionResponse, error)
	StopChargingSession(ctx context.Context, sessionID SessionID, params *StopChargingSessionParams, reqEditors ...RequestEditorFn) (*StopChargingSessionResponse, error)
	GetCurrentDriver(ctx context.Context, reqEditors ...RequestEditorFn) (*GetCurrentDriverResponse, error)
	ListFleetInvoices(ctx context.Context, fleetID FleetID, params *ListFleetInvoicesParams, reqEditors ...RequestEditorFn) (*ListFleetInvoicesResponse, error)
	ListFleetVehicles(ctx context.Context, fleetID FleetID, params *ListFleetVehiclesParams, reqEditors ...RequestEditorFn) (*ListFleetVehiclesResponse, error)
	GetHealth(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHealthResponse, error)
	SearchChargingStations(ctx context.Context, params *SearchChargingStationsParams, reqEditors ...RequestEditorFn) (*SearchChargingStationsResponse, error)
}

type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type rawClientInterface interface {
	startChargingSession(ctx context.Context, params *StartChargingSessionParams, body StartChargingSessionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
	getChargingSession(ctx context.Context, sessionID SessionID, reqEditors ...RequestEditorFn) (*http.Response, error)
	stopChargingSession(ctx context.Context, sessionID SessionID, params *StopChargingSessionParams, reqEditors ...RequestEditorFn) (*http.Response, error)
	getCurrentDriver(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
	listFleetInvoices(ctx context.Context, fleetID FleetID, params *ListFleetInvoicesParams, reqEditors ...RequestEditorFn) (*http.Response, error)
	listFleetVehicles(ctx context.Context, fleetID FleetID, params *ListFleetVehiclesParams, reqEditors ...RequestEditorFn) (*http.Response, error)
	getHealth(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
	searchChargingStations(ctx context.Context, params *SearchChargingStationsParams, reqEditors ...RequestEditorFn) (*http.Response, error)
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

func (c *Client) GetChargingSession(ctx context.Context, sessionID SessionID, reqEditors ...RequestEditorFn) (*GetChargingSessionResponse, error) {
	rsp, err := c.getChargingSession(ctx, sessionID, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseGetChargingSessionResponse(rsp)
}

func (c *Client) GetCurrentDriver(ctx context.Context, reqEditors ...RequestEditorFn) (*GetCurrentDriverResponse, error) {
	rsp, err := c.getCurrentDriver(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseGetCurrentDriverResponse(rsp)
}

func (c *Client) GetHealth(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHealthResponse, error) {
	rsp, err := c.getHealth(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseGetHealthResponse(rsp)
}

func (c *Client) ListFleetInvoices(ctx context.Context, fleetID FleetID, params *ListFleetInvoicesParams, reqEditors ...RequestEditorFn) (*ListFleetInvoicesResponse, error) {
	rsp, err := c.listFleetInvoices(ctx, fleetID, params, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseListFleetInvoicesResponse(rsp)
}

func (c *Client) ListFleetVehicles(ctx context.Context, fleetID FleetID, params *ListFleetVehiclesParams, reqEditors ...RequestEditorFn) (*ListFleetVehiclesResponse, error) {
	rsp, err := c.listFleetVehicles(ctx, fleetID, params, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseListFleetVehiclesResponse(rsp)
}

func (c *Client) SearchChargingStations(ctx context.Context, params *SearchChargingStationsParams, reqEditors ...RequestEditorFn) (*SearchChargingStationsResponse, error) {
	rsp, err := c.searchChargingStations(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseSearchChargingStationsResponse(rsp)
}

func (c *Client) StartChargingSession(ctx context.Context, params *StartChargingSessionParams, body StartChargingSessionJSONRequestBody, reqEditors ...RequestEditorFn) (*StartChargingSessionResponse, error) {
	rsp, err := c.startChargingSession(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseStartChargingSessionResponse(rsp)
}

func (c *Client) StopChargingSession(ctx context.Context, sessionID SessionID, params *StopChargingSessionParams, reqEditors ...RequestEditorFn) (*StopChargingSessionResponse, error) {
	rsp, err := c.stopChargingSession(ctx, sessionID, params, reqEditors...)
	if err != nil {
		return nil, err
	}

	return parseStopChargingSessionResponse(rsp)
}

type ClientOption func(*rawClient) error

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

func WithDriverOAuth(clientID, clientSecret, username, password string, oidcIssuerURL ...string) ClientOption {
	return func(c *rawClient) error {
		ctx := context.Background()
		config := &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint: oauth2.Endpoint{
				TokenURL: "https://identity.voltgrid.example/oauth2/token",
			},
			Scopes: []string{
				"charging:read",
				"charging:start",
				"charging:stop",
				"profile:read",
				"stations:read",
			},
		}

		if len(oidcIssuerURL) > 1 {
			return fmt.Errorf("expected at most one OIDC issuer URL")
		}

		if len(oidcIssuerURL) == 1 && oidcIssuerURL[0] != "" {
			provider, err := oidc.NewProvider(ctx, oidcIssuerURL[0])
			if err != nil {
				return fmt.Errorf("discover OIDC provider: %w", err)
			}

			endpoint := provider.Endpoint()
			if endpoint.TokenURL == "" {
				return fmt.Errorf("OIDC provider has no token endpoint")
			}

			config.Endpoint = endpoint
		}

		token, err := config.PasswordCredentialsToken(ctx, username, password)
		if err != nil {
			return fmt.Errorf("authenticate: %w", err)
		}

		c.DriverOAuth.TokenSource = config.TokenSource(ctx, token)

		return nil
	}
}

func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *rawClient) error {
		c.Client = doer
		return nil
	}
}

func WithPartnerOAuth(clientID, clientSecret string, oidcIssuerURL ...string) ClientOption {
	return func(c *rawClient) error {
		ctx := context.Background()
		config := &clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			TokenURL:     "https://identity.voltgrid.example/oauth2/token",
			Scopes: []string{
				"fleet:billing:read",
				"fleet:charging:read",
				"fleet:charging:start",
				"fleet:charging:stop",
				"fleet:vehicles:read",
				"stations:read",
			},
		}

		if len(oidcIssuerURL) > 1 {
			return fmt.Errorf("expected at most one OIDC issuer URL")
		}

		if len(oidcIssuerURL) == 1 && oidcIssuerURL[0] != "" {
			provider, err := oidc.NewProvider(ctx, oidcIssuerURL[0])
			if err != nil {
				return fmt.Errorf("discover OIDC provider: %w", err)
			}

			endpoint := provider.Endpoint()
			if endpoint.TokenURL == "" {
				return fmt.Errorf("OIDC provider has no token endpoint")
			}

			config.TokenURL = endpoint.TokenURL
			config.AuthStyle = endpoint.AuthStyle
		}

		c.PartnerOAuth.TokenSource = config.TokenSource(ctx)

		return nil
	}
}

func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *rawClient) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

type GetChargingSessionResponse struct {
	Body                      []byte
	HTTPResponse              *http.Response
	JSON200                   *ChargingSession
	ApplicationProblemJSON401 *Unauthorized
	ApplicationProblemJSON403 *Forbidden
	ApplicationProblemJSON404 *NotFound
	Headers401                *GetChargingSessionResponse401Headers
}

func parseGetChargingSessionResponse(rsp *http.Response) (*GetChargingSessionResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &GetChargingSessionResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}
	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest ChargingSession
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON200 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON401 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON403 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON404 = &dest
	}

	switch {
	case rsp.StatusCode == 401:
		var headers GetChargingSessionResponse401Headers

		if values := rsp.Header.Values("WWW-Authenticate"); len(values) > 0 {
			var value string
			if err := runtime.BindStyledParameterWithOptions("simple", "WWW-Authenticate", values[0], &value, runtime.BindStyledParameterOptions{
				ParamLocation: runtime.ParamLocationHeader,
				Explode:       false,
				Required:      false,
				Type:          "string",
				Format:        "",
			}); err != nil {
				return nil, err
			}

			headers.WWWAuthenticate = &value
		}

		response.Headers401 = &headers
	}

	return response, nil
}

func (r GetChargingSessionResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

func (r GetChargingSessionResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r GetChargingSessionResponse) GetBody() []byte {
	return r.Body
}

func (r GetChargingSessionResponse) GetChargingSession() (*ChargingSession, error) {
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

func (r GetChargingSessionResponse) GetForbidden() (*Forbidden, error) {
	if r.ApplicationProblemJSON403 != nil {
		return r.ApplicationProblemJSON403, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 403:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r GetChargingSessionResponse) GetNotFound() (*NotFound, error) {
	if r.ApplicationProblemJSON404 != nil {
		return r.ApplicationProblemJSON404, nil
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

func (r GetChargingSessionResponse) GetUnauthorized() (*Unauthorized, error) {
	if r.ApplicationProblemJSON401 != nil {
		return r.ApplicationProblemJSON401, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 401:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r GetChargingSessionResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r GetChargingSessionResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

type GetChargingSessionResponse401Headers struct {
	WWWAuthenticate *string
}

type GetCurrentDriverResponse struct {
	Body                      []byte
	HTTPResponse              *http.Response
	JSON200                   *Driver
	ApplicationProblemJSON401 *Unauthorized
	ApplicationProblemJSON403 *Forbidden
	Headers401                *GetCurrentDriverResponse401Headers
}

func parseGetCurrentDriverResponse(rsp *http.Response) (*GetCurrentDriverResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &GetCurrentDriverResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}
	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Driver
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON200 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON401 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON403 = &dest
	}

	switch {
	case rsp.StatusCode == 401:
		var headers GetCurrentDriverResponse401Headers

		if values := rsp.Header.Values("WWW-Authenticate"); len(values) > 0 {
			var value string
			if err := runtime.BindStyledParameterWithOptions("simple", "WWW-Authenticate", values[0], &value, runtime.BindStyledParameterOptions{
				ParamLocation: runtime.ParamLocationHeader,
				Explode:       false,
				Required:      false,
				Type:          "string",
				Format:        "",
			}); err != nil {
				return nil, err
			}

			headers.WWWAuthenticate = &value
		}

		response.Headers401 = &headers
	}

	return response, nil
}

func (r GetCurrentDriverResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

func (r GetCurrentDriverResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r GetCurrentDriverResponse) GetBody() []byte {
	return r.Body
}

func (r GetCurrentDriverResponse) GetDriver() (*Driver, error) {
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

func (r GetCurrentDriverResponse) GetForbidden() (*Forbidden, error) {
	if r.ApplicationProblemJSON403 != nil {
		return r.ApplicationProblemJSON403, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 403:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r GetCurrentDriverResponse) GetUnauthorized() (*Unauthorized, error) {
	if r.ApplicationProblemJSON401 != nil {
		return r.ApplicationProblemJSON401, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 401:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r GetCurrentDriverResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r GetCurrentDriverResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

type GetCurrentDriverResponse401Headers struct {
	WWWAuthenticate *string
}

type GetHealthResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Status GetHealth200JSONResponseBodyStatus `json:"status"`
	}
}

func parseGetHealthResponse(rsp *http.Response) (*GetHealthResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &GetHealthResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}
	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Status GetHealth200JSONResponseBodyStatus `json:"status"`
		}
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON200 = &dest
	}

	return response, nil
}

func (r GetHealthResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

func (r GetHealthResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r GetHealthResponse) GetBody() []byte {
	return r.Body
}

func (r GetHealthResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r GetHealthResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

type ListFleetInvoicesResponse struct {
	Body                      []byte
	HTTPResponse              *http.Response
	JSON200                   *[]Invoice
	ApplicationProblemJSON401 *Unauthorized
	ApplicationProblemJSON403 *Forbidden
	Headers401                *ListFleetInvoicesResponse401Headers
}

func parseListFleetInvoicesResponse(rsp *http.Response) (*ListFleetInvoicesResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &ListFleetInvoicesResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}
	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Invoice
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON200 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON401 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON403 = &dest
	}

	switch {
	case rsp.StatusCode == 401:
		var headers ListFleetInvoicesResponse401Headers

		if values := rsp.Header.Values("WWW-Authenticate"); len(values) > 0 {
			var value string
			if err := runtime.BindStyledParameterWithOptions("simple", "WWW-Authenticate", values[0], &value, runtime.BindStyledParameterOptions{
				ParamLocation: runtime.ParamLocationHeader,
				Explode:       false,
				Required:      false,
				Type:          "string",
				Format:        "",
			}); err != nil {
				return nil, err
			}

			headers.WWWAuthenticate = &value
		}

		response.Headers401 = &headers
	}

	return response, nil
}

func (r ListFleetInvoicesResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

func (r ListFleetInvoicesResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r ListFleetInvoicesResponse) GetBody() []byte {
	return r.Body
}

func (r ListFleetInvoicesResponse) GetForbidden() (*Forbidden, error) {
	if r.ApplicationProblemJSON403 != nil {
		return r.ApplicationProblemJSON403, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 403:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r ListFleetInvoicesResponse) GetInvoices() (*[]Invoice, error) {
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

func (r ListFleetInvoicesResponse) GetUnauthorized() (*Unauthorized, error) {
	if r.ApplicationProblemJSON401 != nil {
		return r.ApplicationProblemJSON401, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 401:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r ListFleetInvoicesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r ListFleetInvoicesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

type ListFleetInvoicesResponse401Headers struct {
	WWWAuthenticate *string
}

type ListFleetVehiclesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Items      []Vehicle `json:"items"`
		NextCursor *string   `json:"nextCursor"`
	}
	ApplicationProblemJSON401 *Unauthorized
	ApplicationProblemJSON403 *Forbidden
	ApplicationProblemJSON404 *NotFound
	Headers401                *ListFleetVehiclesResponse401Headers
}

func parseListFleetVehiclesResponse(rsp *http.Response) (*ListFleetVehiclesResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &ListFleetVehiclesResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}
	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Items      []Vehicle `json:"items"`
			NextCursor *string   `json:"nextCursor"`
		}
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON200 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON401 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON403 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON404 = &dest
	}

	switch {
	case rsp.StatusCode == 401:
		var headers ListFleetVehiclesResponse401Headers

		if values := rsp.Header.Values("WWW-Authenticate"); len(values) > 0 {
			var value string
			if err := runtime.BindStyledParameterWithOptions("simple", "WWW-Authenticate", values[0], &value, runtime.BindStyledParameterOptions{
				ParamLocation: runtime.ParamLocationHeader,
				Explode:       false,
				Required:      false,
				Type:          "string",
				Format:        "",
			}); err != nil {
				return nil, err
			}

			headers.WWWAuthenticate = &value
		}

		response.Headers401 = &headers
	}

	return response, nil
}

func (r ListFleetVehiclesResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

func (r ListFleetVehiclesResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r ListFleetVehiclesResponse) GetBody() []byte {
	return r.Body
}

func (r ListFleetVehiclesResponse) GetForbidden() (*Forbidden, error) {
	if r.ApplicationProblemJSON403 != nil {
		return r.ApplicationProblemJSON403, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 403:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r ListFleetVehiclesResponse) GetNotFound() (*NotFound, error) {
	if r.ApplicationProblemJSON404 != nil {
		return r.ApplicationProblemJSON404, nil
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

func (r ListFleetVehiclesResponse) GetUnauthorized() (*Unauthorized, error) {
	if r.ApplicationProblemJSON401 != nil {
		return r.ApplicationProblemJSON401, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 401:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r ListFleetVehiclesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r ListFleetVehiclesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

type ListFleetVehiclesResponse401Headers struct {
	WWWAuthenticate *string
}

type RequestEditorFn func(ctx context.Context, req *http.Request) error

type SearchChargingStationsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Items []Station `json:"items"`
		Total int       `json:"total"`
	}
	ApplicationProblemJSON401 *Unauthorized
	ApplicationProblemJSON403 *Forbidden
	ApplicationProblemJSON429 *RateLimited
	Headers401                *SearchChargingStationsResponse401Headers
	Headers429                *SearchChargingStationsResponse429Headers
}

func parseSearchChargingStationsResponse(rsp *http.Response) (*SearchChargingStationsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &SearchChargingStationsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}
	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Items []Station `json:"items"`
			Total int       `json:"total"`
		}
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON200 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON401 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON403 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest RateLimited
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON429 = &dest
	}

	switch {
	case rsp.StatusCode == 401:
		var headers SearchChargingStationsResponse401Headers

		if values := rsp.Header.Values("WWW-Authenticate"); len(values) > 0 {
			var value string
			if err := runtime.BindStyledParameterWithOptions("simple", "WWW-Authenticate", values[0], &value, runtime.BindStyledParameterOptions{
				ParamLocation: runtime.ParamLocationHeader,
				Explode:       false,
				Required:      false,
				Type:          "string",
				Format:        "",
			}); err != nil {
				return nil, err
			}

			headers.WWWAuthenticate = &value
		}

		response.Headers401 = &headers
	case rsp.StatusCode == 429:
		var headers SearchChargingStationsResponse429Headers

		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{
				ParamLocation: runtime.ParamLocationHeader,
				Explode:       false,
				Required:      false,
				Type:          "integer",
				Format:        "",
			}); err != nil {
				return nil, err
			}

			headers.RetryAfter = &value
		}

		response.Headers429 = &headers
	}

	return response, nil
}

func (r SearchChargingStationsResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

func (r SearchChargingStationsResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r SearchChargingStationsResponse) GetBody() []byte {
	return r.Body
}

func (r SearchChargingStationsResponse) GetForbidden() (*Forbidden, error) {
	if r.ApplicationProblemJSON403 != nil {
		return r.ApplicationProblemJSON403, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 403:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r SearchChargingStationsResponse) GetRateLimited() (*RateLimited, error) {
	if r.ApplicationProblemJSON429 != nil {
		return r.ApplicationProblemJSON429, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 429:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r SearchChargingStationsResponse) GetUnauthorized() (*Unauthorized, error) {
	if r.ApplicationProblemJSON401 != nil {
		return r.ApplicationProblemJSON401, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 401:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r SearchChargingStationsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r SearchChargingStationsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

type SearchChargingStationsResponse401Headers struct {
	WWWAuthenticate *string
}

type SearchChargingStationsResponse429Headers struct {
	RetryAfter *int
}

type StartChargingSessionResponse struct {
	Body                      []byte
	HTTPResponse              *http.Response
	JSON201                   *ChargingSession
	ApplicationProblemJSON400 *BadRequest
	ApplicationProblemJSON401 *Unauthorized
	ApplicationProblemJSON403 *Forbidden
	ApplicationProblemJSON409 *Problem
	Headers201                *StartChargingSessionResponse201Headers
	Headers401                *StartChargingSessionResponse401Headers
}

func parseStartChargingSessionResponse(rsp *http.Response) (*StartChargingSessionResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &StartChargingSessionResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}
	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest ChargingSession
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON201 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON400 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON401 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON403 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 409:
		var dest Problem
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON409 = &dest
	}

	switch {
	case rsp.StatusCode == 201:
		var headers StartChargingSessionResponse201Headers

		if values := rsp.Header.Values("Location"); len(values) > 0 {
			var value string
			if err := runtime.BindStyledParameterWithOptions("simple", "Location", values[0], &value, runtime.BindStyledParameterOptions{
				ParamLocation: runtime.ParamLocationHeader,
				Explode: false,
				Required: false,
				Type: "string",
				Format: "uri-reference",
			}); err != nil {
				return nil, err
			}

			headers.Location = &value
		}

		response.Headers201 = &headers
	case rsp.StatusCode == 401:
		var headers StartChargingSessionResponse401Headers

		if values := rsp.Header.Values("WWW-Authenticate"); len(values) > 0 {
			var value string
			if err := runtime.BindStyledParameterWithOptions("simple", "WWW-Authenticate", values[0], &value, runtime.BindStyledParameterOptions{
				ParamLocation: runtime.ParamLocationHeader,
				Explode: false,
				Required: false,
				Type: "string",
				Format: "",
			}); err != nil {
				return nil, err
			}

			headers.WWWAuthenticate = &value
		}

		response.Headers401 = &headers
	}

	return response, nil
}

func (r StartChargingSessionResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

func (r StartChargingSessionResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r StartChargingSessionResponse) GetBadRequest() (*BadRequest, error) {
	if r.ApplicationProblemJSON400 != nil {
		return r.ApplicationProblemJSON400, nil
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

func (r StartChargingSessionResponse) GetBody() []byte {
	return r.Body
}

func (r StartChargingSessionResponse) GetChargingSession() (*ChargingSession, error) {
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

func (r StartChargingSessionResponse) GetForbidden() (*Forbidden, error) {
	if r.ApplicationProblemJSON403 != nil {
		return r.ApplicationProblemJSON403, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 403:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r StartChargingSessionResponse) GetProblem() (*Problem, error) {
	if r.ApplicationProblemJSON409 != nil {
		return r.ApplicationProblemJSON409, nil
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

func (r StartChargingSessionResponse) GetUnauthorized() (*Unauthorized, error) {
	if r.ApplicationProblemJSON401 != nil {
		return r.ApplicationProblemJSON401, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 401:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r StartChargingSessionResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r StartChargingSessionResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

type StartChargingSessionResponse201Headers struct {
	Location *string
}

type StartChargingSessionResponse401Headers struct {
	WWWAuthenticate *string
}

type StopChargingSessionResponse struct {
	Body                      []byte
	HTTPResponse              *http.Response
	JSON202                   *ChargingSession
	ApplicationProblemJSON401 *Unauthorized
	ApplicationProblemJSON403 *Forbidden
	ApplicationProblemJSON404 *NotFound
	Headers401                *StopChargingSessionResponse401Headers
}

func parseStopChargingSessionResponse(rsp *http.Response) (*StopChargingSessionResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	response := &StopChargingSessionResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}
	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 202:
		var dest ChargingSession
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.JSON202 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON401 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON403 = &dest
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := runtime.JSONUnmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}

		response.ApplicationProblemJSON404 = &dest
	}

	switch {
	case rsp.StatusCode == 401:
		var headers StopChargingSessionResponse401Headers

		if values := rsp.Header.Values("WWW-Authenticate"); len(values) > 0 {
			var value string
			if err := runtime.BindStyledParameterWithOptions("simple", "WWW-Authenticate", values[0], &value, runtime.BindStyledParameterOptions{
				ParamLocation: runtime.ParamLocationHeader,
				Explode: false,
				Required: false,
				Type: "string",
				Format: "",
			}); err != nil {
				return nil, err
			}

			headers.WWWAuthenticate = &value
		}

		response.Headers401 = &headers
	}

	return response, nil
}

func (r StopChargingSessionResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}

	return ""
}

func (r StopChargingSessionResponse) Error() error {
	if r.HTTPResponse == nil {
		return fmt.Errorf("missing HTTP response")
	}

	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
}

func (r StopChargingSessionResponse) GetBody() []byte {
	return r.Body
}

func (r StopChargingSessionResponse) GetChargingSession() (*ChargingSession, error) {
	if r.JSON202 != nil {
		return r.JSON202, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 202:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r StopChargingSessionResponse) GetForbidden() (*Forbidden, error) {
	if r.ApplicationProblemJSON403 != nil {
		return r.ApplicationProblemJSON403, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 403:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r StopChargingSessionResponse) GetNotFound() (*NotFound, error) {
	if r.ApplicationProblemJSON404 != nil {
		return r.ApplicationProblemJSON404, nil
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

func (r StopChargingSessionResponse) GetUnauthorized() (*Unauthorized, error) {
	if r.ApplicationProblemJSON401 != nil {
		return r.ApplicationProblemJSON401, nil
	}

	if r.HTTPResponse == nil {
		return nil, fmt.Errorf("missing HTTP response")
	}

	switch r.HTTPResponse.StatusCode {
	case 401:
		return nil, fmt.Errorf("expected response body for status %d, but no matching response body was found", r.HTTPResponse.StatusCode)
	default:
		return nil, fmt.Errorf("unexpected response: status %d: %s", r.HTTPResponse.StatusCode, string(r.Body))
	}
}

func (r StopChargingSessionResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}

	return http.StatusText(0)
}

func (r StopChargingSessionResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}

	return 0
}

type StopChargingSessionResponse401Headers struct {
	WWWAuthenticate *string
}

type rawClient struct {
	Server         string
	Client         HttpRequestDoer
	RequestEditors []RequestEditorFn
	DriverOAuth    struct {
		TokenSource oauth2.TokenSource
	}
	PartnerOAuth struct {
		TokenSource oauth2.TokenSource
	}
}

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

func (c *rawClient) getChargingSession(ctx context.Context, sessionID SessionID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	editors := append(slices.Clone(reqEditors), func(ctx context.Context, req *http.Request) error {
		switch {
		case c.DriverOAuth.TokenSource != nil:
			token, err := c.DriverOAuth.TokenSource.Token()
			if err != nil {
				return fmt.Errorf("get authentication token for GET /charging-sessions/{sessionId}: %w", err)
			}

			token.SetAuthHeader(req)
		case c.PartnerOAuth.TokenSource != nil:
			token, err := c.PartnerOAuth.TokenSource.Token()
			if err != nil {
				return fmt.Errorf("get authentication token for GET /charging-sessions/{sessionId}: %w", err)
			}

			token.SetAuthHeader(req)
		}

		return nil
	})

	return c.unauthenticatedGetChargingSession(ctx, sessionID, editors...)
}

func (c *rawClient) getCurrentDriver(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	editors := append(slices.Clone(reqEditors), func(ctx context.Context, req *http.Request) error {
		switch {
		case c.DriverOAuth.TokenSource != nil:
			token, err := c.DriverOAuth.TokenSource.Token()
			if err != nil {
				return fmt.Errorf("get authentication token for GET /drivers/me: %w", err)
			}

			token.SetAuthHeader(req)
		}

		return nil
	})

	return c.unauthenticatedGetCurrentDriver(ctx, editors...)
}

func (c *rawClient) getHealth(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newGetHealthRequest(c.Server)
	if err != nil {
		return nil, fmt.Errorf("create GET /health request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /health request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /health request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) listFleetInvoices(ctx context.Context, fleetID FleetID, params *ListFleetInvoicesParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	editors := append(slices.Clone(reqEditors), func(ctx context.Context, req *http.Request) error {
		switch {
		case c.PartnerOAuth.TokenSource != nil:
			token, err := c.PartnerOAuth.TokenSource.Token()
			if err != nil {
				return fmt.Errorf("get authentication token for GET /fleets/{fleetId}/invoices: %w", err)
			}

			token.SetAuthHeader(req)
		}

		return nil
	})

	return c.unauthenticatedListFleetInvoices(ctx, fleetID, params, editors...)
}

func (c *rawClient) listFleetVehicles(ctx context.Context, fleetID FleetID, params *ListFleetVehiclesParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	editors := append(slices.Clone(reqEditors), func(ctx context.Context, req *http.Request) error {
		switch {
		case c.PartnerOAuth.TokenSource != nil:
			token, err := c.PartnerOAuth.TokenSource.Token()
			if err != nil {
				return fmt.Errorf("get authentication token for GET /fleets/{fleetId}/vehicles: %w", err)
			}

			token.SetAuthHeader(req)
		}

		return nil
	})

	return c.unauthenticatedListFleetVehicles(ctx, fleetID, params, editors...)
}

func (c *rawClient) searchChargingStations(ctx context.Context, params *SearchChargingStationsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	editors := append(slices.Clone(reqEditors), func(ctx context.Context, req *http.Request) error {
		switch {
		case c.DriverOAuth.TokenSource != nil:
			token, err := c.DriverOAuth.TokenSource.Token()
			if err != nil {
				return fmt.Errorf("get authentication token for GET /stations: %w", err)
			}

			token.SetAuthHeader(req)
		case c.PartnerOAuth.TokenSource != nil:
			token, err := c.PartnerOAuth.TokenSource.Token()
			if err != nil {
				return fmt.Errorf("get authentication token for GET /stations: %w", err)
			}

			token.SetAuthHeader(req)
		}

		return nil
	})

	return c.unauthenticatedSearchChargingStations(ctx, params, editors...)
}

func (c *rawClient) startChargingSession(ctx context.Context, params *StartChargingSessionParams, body StartChargingSessionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	editors := append(slices.Clone(reqEditors), func(ctx context.Context, req *http.Request) error {
		switch {
		case c.DriverOAuth.TokenSource != nil:
			token, err := c.DriverOAuth.TokenSource.Token()
			if err != nil {
				return fmt.Errorf("get authentication token for POST /charging-sessions: %w", err)
			}

			token.SetAuthHeader(req)
		case c.PartnerOAuth.TokenSource != nil:
			token, err := c.PartnerOAuth.TokenSource.Token()
			if err != nil {
				return fmt.Errorf("get authentication token for POST /charging-sessions: %w", err)
			}

			token.SetAuthHeader(req)
		}

		return nil
	})

	return c.unauthenticatedStartChargingSession(ctx, params, body, editors...)
}

func (c *rawClient) stopChargingSession(ctx context.Context, sessionID SessionID, params *StopChargingSessionParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	editors := append(slices.Clone(reqEditors), func(ctx context.Context, req *http.Request) error {
		switch {
		case c.DriverOAuth.TokenSource != nil:
			token, err := c.DriverOAuth.TokenSource.Token()
			if err != nil {
				return fmt.Errorf("get authentication token for POST /charging-sessions/{sessionId}/stop: %w", err)
			}

			token.SetAuthHeader(req)
		case c.PartnerOAuth.TokenSource != nil:
			token, err := c.PartnerOAuth.TokenSource.Token()
			if err != nil {
				return fmt.Errorf("get authentication token for POST /charging-sessions/{sessionId}/stop: %w", err)
			}

			token.SetAuthHeader(req)
		}

		return nil
	})

	return c.unauthenticatedStopChargingSession(ctx, sessionID, params, editors...)
}

func (c *rawClient) unauthenticatedGetChargingSession(ctx context.Context, sessionID SessionID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newGetChargingSessionRequest(c.Server, sessionID)
	if err != nil {
		return nil, fmt.Errorf("create GET /charging-sessions/{sessionId} request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /charging-sessions/{sessionId} request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /charging-sessions/{sessionId} request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) unauthenticatedGetCurrentDriver(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newGetCurrentDriverRequest(c.Server)
	if err != nil {
		return nil, fmt.Errorf("create GET /drivers/me request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /drivers/me request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /drivers/me request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) unauthenticatedListFleetInvoices(ctx context.Context, fleetID FleetID, params *ListFleetInvoicesParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newListFleetInvoicesRequest(c.Server, fleetID, params)
	if err != nil {
		return nil, fmt.Errorf("create GET /fleets/{fleetId}/invoices request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /fleets/{fleetId}/invoices request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /fleets/{fleetId}/invoices request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) unauthenticatedListFleetVehicles(ctx context.Context, fleetID FleetID, params *ListFleetVehiclesParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newListFleetVehiclesRequest(c.Server, fleetID, params)
	if err != nil {
		return nil, fmt.Errorf("create GET /fleets/{fleetId}/vehicles request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /fleets/{fleetId}/vehicles request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /fleets/{fleetId}/vehicles request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) unauthenticatedSearchChargingStations(ctx context.Context, params *SearchChargingStationsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newSearchChargingStationsRequest(c.Server, params)
	if err != nil {
		return nil, fmt.Errorf("create GET /stations request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to GET /stations request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform GET /stations request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) unauthenticatedStartChargingSession(ctx context.Context, params *StartChargingSessionParams, body StartChargingSessionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newStartChargingSessionRequest(c.Server, params, body)
	if err != nil {
		return nil, fmt.Errorf("create POST /charging-sessions request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to POST /charging-sessions request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform POST /charging-sessions request: %w", err)
	}

	return resp, nil
}

func (c *rawClient) unauthenticatedStopChargingSession(ctx context.Context, sessionID SessionID, params *StopChargingSessionParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newStopChargingSessionRequest(c.Server, sessionID, params)
	if err != nil {
		return nil, fmt.Errorf("create POST /charging-sessions/{sessionId}/stop request: %w", err)
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, fmt.Errorf("apply editors to POST /charging-sessions/{sessionId}/stop request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform POST /charging-sessions/{sessionId}/stop request: %w", err)
	}

	return resp, nil
}

func newGetChargingSessionRequest(server string, sessionID SessionID) (*http.Request, error) {
	var err error
	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "sessionId", sessionID, runtime.StyleParamOptions{
		ParamLocation: runtime.ParamLocationPath,
		Type:          "string",
		Format:        "",
	})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/charging-sessions/" + pathParam0
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

func newGetCurrentDriverRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/drivers/me"
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

func newListFleetInvoicesRequest(server string, fleetID FleetID, params *ListFleetInvoicesParams) (*http.Request, error) {
	var err error
	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "fleetId", fleetID, runtime.StyleParamOptions{
		ParamLocation: runtime.ParamLocationPath,
		Type:          "string",
		Format:        "",
	})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/fleets/" + pathParam0 + "/invoices"
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

		if queryFrag, err := runtime.StyleParamWithOptions("form", true, "billingPeriod", params.BillingPeriod, runtime.StyleParamOptions{
			ParamLocation: runtime.ParamLocationQuery,
			Type:          "string",
			Format:        "",
		}); err != nil {
			return nil, err
		} else {
			for _, qp := range strings.Split(queryFrag, "&") {
				rawQueryFragments = append(rawQueryFragments, qp)
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

func newListFleetVehiclesRequest(server string, fleetID FleetID, params *ListFleetVehiclesParams) (*http.Request, error) {
	var err error
	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "fleetId", fleetID, runtime.StyleParamOptions{
		ParamLocation: runtime.ParamLocationPath,
		Type:          "string",
		Format:        "",
	})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/fleets/" + pathParam0 + "/vehicles"
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

		if params.Limit != nil {
			if queryFrag, err := runtime.StyleParamWithOptions("form", true, "limit", *params.Limit, runtime.StyleParamOptions{
				ParamLocation: runtime.ParamLocationQuery,
				Type:          "integer",
				Format:        "",
			}); err != nil {
				return nil, err
			} else {
				for _, qp := range strings.Split(queryFrag, "&") {
					rawQueryFragments = append(rawQueryFragments, qp)
				}
			}
		}

		if params.Cursor != nil {
			if queryFrag, err := runtime.StyleParamWithOptions("form", true, "cursor", *params.Cursor, runtime.StyleParamOptions{
				ParamLocation: runtime.ParamLocationQuery,
				Type:          "string",
				Format:        "",
			}); err != nil {
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

func newSearchChargingStationsRequest(server string, params *SearchChargingStationsParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/stations"
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

		if queryFrag, err := runtime.StyleParamWithOptions("form", true, "latitude", params.Latitude, runtime.StyleParamOptions{
			ParamLocation: runtime.ParamLocationQuery,
			Type:          "number",
			Format:        "",
		}); err != nil {
			return nil, err
		} else {
			for _, qp := range strings.Split(queryFrag, "&") {
				rawQueryFragments = append(rawQueryFragments, qp)
			}
		}

		if queryFrag, err := runtime.StyleParamWithOptions("form", true, "longitude", params.Longitude, runtime.StyleParamOptions{
			ParamLocation: runtime.ParamLocationQuery,
			Type:          "number",
			Format:        "",
		}); err != nil {
			return nil, err
		} else {
			for _, qp := range strings.Split(queryFrag, "&") {
				rawQueryFragments = append(rawQueryFragments, qp)
			}
		}

		if params.RadiusKm != nil {
			if queryFrag, err := runtime.StyleParamWithOptions("form", true, "radiusKm", *params.RadiusKm, runtime.StyleParamOptions{
				ParamLocation: runtime.ParamLocationQuery,
				Type:          "number",
				Format:        "",
			}); err != nil {
				return nil, err
			} else {
				for _, qp := range strings.Split(queryFrag, "&") {
					rawQueryFragments = append(rawQueryFragments, qp)
				}
			}
		}

		if params.AvailableOnly != nil {
			if queryFrag, err := runtime.StyleParamWithOptions("form", true, "availableOnly", *params.AvailableOnly, runtime.StyleParamOptions{
				ParamLocation: runtime.ParamLocationQuery,
				Type:          "boolean",
				Format:        "",
			}); err != nil {
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

func newStartChargingSessionRequest(server string, params *StartChargingSessionParams, body StartChargingSessionJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader

	buf, err := runtime.JSONMarshal(body)
	if err != nil {
		return nil, err
	}

	bodyReader = bytes.NewReader(buf)

	return newStartChargingSessionRequestWithBody(server, params, "application/json", bodyReader)
}

func newStartChargingSessionRequestWithBody(server string, params *StartChargingSessionParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/charging-sessions"
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

	if params != nil {
		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithOptions("simple", false, "Idempotency-Key", params.IdempotencyKey, runtime.StyleParamOptions{
			ParamLocation: runtime.ParamLocationHeader,
			Type:          "string",
			Format:        "uuid",
		})
		if err != nil {
			return nil, err
		}

		req.Header.Set("Idempotency-Key", headerParam0)
	}

	return req, nil
}

func newStopChargingSessionRequest(server string, sessionID SessionID, params *StopChargingSessionParams) (*http.Request, error) {
	var err error
	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "sessionId", sessionID, runtime.StyleParamOptions{
		ParamLocation: runtime.ParamLocationPath,
		Type:          "string",
		Format:        "",
	})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/charging-sessions/" + pathParam0 + "/stop"
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

	if params != nil {
		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithOptions("simple", false, "Idempotency-Key", params.IdempotencyKey, runtime.StyleParamOptions{
			ParamLocation: runtime.ParamLocationHeader,
			Type:          "string",
			Format:        "uuid",
		})
		if err != nil {
			return nil, err
		}

		req.Header.Set("Idempotency-Key", headerParam0)
	}

	return req, nil
}
