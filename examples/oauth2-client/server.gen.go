package examples

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	chi "github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type ServerInterface interface {
	StartChargingSession(w http.ResponseWriter, r *http.Request, params StartChargingSessionParams)
	GetChargingSession(w http.ResponseWriter, r *http.Request, sessionID SessionID)
	StopChargingSession(w http.ResponseWriter, r *http.Request, sessionID SessionID, params StopChargingSessionParams)
	GetCurrentDriver(w http.ResponseWriter, r *http.Request)
	ListFleetInvoices(w http.ResponseWriter, r *http.Request, fleetID FleetID, params ListFleetInvoicesParams)
	ListFleetVehicles(w http.ResponseWriter, r *http.Request, fleetID FleetID, params ListFleetVehiclesParams)
	GetHealth(w http.ResponseWriter, r *http.Request)
	SearchChargingStations(w http.ResponseWriter, r *http.Request, params SearchChargingStationsParams)
}

type Unimplemented struct{}

func (_ Unimplemented) StartChargingSession(w http.ResponseWriter, r *http.Request, params StartChargingSessionParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) GetChargingSession(w http.ResponseWriter, r *http.Request, sessionID SessionID) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) StopChargingSession(w http.ResponseWriter, r *http.Request, sessionID SessionID, params StopChargingSessionParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) GetCurrentDriver(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) ListFleetInvoices(w http.ResponseWriter, r *http.Request, fleetID FleetID, params ListFleetInvoicesParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) ListFleetVehicles(w http.ResponseWriter, r *http.Request, fleetID FleetID, params ListFleetVehiclesParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) GetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) SearchChargingStations(w http.ResponseWriter, r *http.Request, params SearchChargingStationsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

func (siw *ServerInterfaceWrapper) StartChargingSession(w http.ResponseWriter, r *http.Request) {
	var err error
	_ = err
	var params StartChargingSessionParams

	headers := r.Header
	if valueList, found := headers[http.CanonicalHeaderKey("Idempotency-Key")]; found {
		var IdempotencyKey IdempotencyKey

		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "Idempotency-Key", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "Idempotency-Key", valueList[0], &IdempotencyKey, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true, Type: "string", Format: "uuid"})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "Idempotency-Key", Err: err})
			return
		}

		params.IdempotencyKey = IdempotencyKey
	} else {
		err := fmt.Errorf("Header parameter Idempotency-Key is required, but not found")
		siw.ErrorHandlerFunc(w, r, &RequiredHeaderError{ParamName: "Idempotency-Key", Err: err})

		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.StartChargingSession(w, r, params)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetChargingSession(w http.ResponseWriter, r *http.Request) {
	var err error
	_ = err
	var sessionID SessionID

	err = runtime.BindStyledParameterWithOptions("simple", "sessionId", chi.URLParam(r, "sessionId"), &sessionID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "string", Format: "", ValueIsUnescaped: r.URL.RawPath == ""})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "sessionId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetChargingSession(w, r, sessionID)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) StopChargingSession(w http.ResponseWriter, r *http.Request) {
	var err error
	_ = err
	var sessionID SessionID

	err = runtime.BindStyledParameterWithOptions("simple", "sessionId", chi.URLParam(r, "sessionId"), &sessionID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "string", Format: "", ValueIsUnescaped: r.URL.RawPath == ""})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "sessionId", Err: err})
		return
	}
	var params StopChargingSessionParams

	headers := r.Header
	if valueList, found := headers[http.CanonicalHeaderKey("Idempotency-Key")]; found {
		var IdempotencyKey IdempotencyKey

		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "Idempotency-Key", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "Idempotency-Key", valueList[0], &IdempotencyKey, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true, Type: "string", Format: "uuid"})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "Idempotency-Key", Err: err})
			return
		}

		params.IdempotencyKey = IdempotencyKey
	} else {
		err := fmt.Errorf("Header parameter Idempotency-Key is required, but not found")
		siw.ErrorHandlerFunc(w, r, &RequiredHeaderError{ParamName: "Idempotency-Key", Err: err})

		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.StopChargingSession(w, r, sessionID, params)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetCurrentDriver(w http.ResponseWriter, r *http.Request) {
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetCurrentDriver(w, r)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) ListFleetInvoices(w http.ResponseWriter, r *http.Request) {
	var err error
	_ = err
	var fleetID FleetID

	err = runtime.BindStyledParameterWithOptions("simple", "fleetId", chi.URLParam(r, "fleetId"), &fleetID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "string", Format: "", ValueIsUnescaped: r.URL.RawPath == ""})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "fleetId", Err: err})
		return
	}
	var params ListFleetInvoicesParams

	err = runtime.BindQueryParameterWithOptions("form", true, true, "billingPeriod", r.URL.Query(), &params.BillingPeriod, runtime.BindQueryParameterOptions{Type: "string", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "billingPeriod"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "billingPeriod", Err: err})
		}

		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListFleetInvoices(w, r, fleetID, params)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) ListFleetVehicles(w http.ResponseWriter, r *http.Request) {
	var err error
	_ = err
	var fleetID FleetID

	err = runtime.BindStyledParameterWithOptions("simple", "fleetId", chi.URLParam(r, "fleetId"), &fleetID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "string", Format: "", ValueIsUnescaped: r.URL.RawPath == ""})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "fleetId", Err: err})
		return
	}
	var params ListFleetVehiclesParams

	err = runtime.BindQueryParameterWithOptions("form", true, false, "limit", r.URL.Query(), &params.Limit, runtime.BindQueryParameterOptions{Type: "integer", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "limit"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		}

		return
	}

	err = runtime.BindQueryParameterWithOptions("form", true, false, "cursor", r.URL.Query(), &params.Cursor, runtime.BindQueryParameterOptions{Type: "string", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "cursor"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "cursor", Err: err})
		}

		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListFleetVehicles(w, r, fleetID, params)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetHealth(w http.ResponseWriter, r *http.Request) {
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetHealth(w, r)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) SearchChargingStations(w http.ResponseWriter, r *http.Request) {
	var err error
	_ = err
	var params SearchChargingStationsParams

	err = runtime.BindQueryParameterWithOptions("form", true, true, "latitude", r.URL.Query(), &params.Latitude, runtime.BindQueryParameterOptions{Type: "number", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "latitude"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "latitude", Err: err})
		}

		return
	}

	err = runtime.BindQueryParameterWithOptions("form", true, true, "longitude", r.URL.Query(), &params.Longitude, runtime.BindQueryParameterOptions{Type: "number", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "longitude"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "longitude", Err: err})
		}

		return
	}

	err = runtime.BindQueryParameterWithOptions("form", true, false, "radiusKm", r.URL.Query(), &params.RadiusKm, runtime.BindQueryParameterOptions{Type: "number", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "radiusKm"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "radiusKm", Err: err})
		}

		return
	}

	err = runtime.BindQueryParameterWithOptions("form", true, false, "availableOnly", r.URL.Query(), &params.AvailableOnly, runtime.BindQueryParameterOptions{Type: "boolean", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "availableOnly"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "availableOnly", Err: err})
		}

		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SearchChargingStations(w, r, params)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter
	if r == nil {
		r = chi.NewRouter()
	}

	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/stations", wrapper.SearchChargingStations)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/drivers/me", wrapper.GetCurrentDriver)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/charging-sessions", wrapper.StartChargingSession)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/charging-sessions/{sessionId}", wrapper.GetChargingSession)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/charging-sessions/{sessionId}/stop", wrapper.StopChargingSession)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/fleets/{fleetId}/vehicles", wrapper.ListFleetVehicles)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/fleets/{fleetId}/invoices", wrapper.ListFleetInvoices)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/health", wrapper.GetHealth)
	})

	return r
}

type BadRequestApplicationProblemPlusJSONResponse Problem

type ForbiddenApplicationProblemPlusJSONResponse Problem

type NotFoundApplicationProblemPlusJSONResponse Problem

type RateLimitedResponseHeaders struct {
	RetryAfter *int
}

type RateLimitedApplicationProblemPlusJSONResponse struct {
	Body    Problem
	Headers RateLimitedResponseHeaders
}

type UnauthorizedResponseHeaders struct {
	WWWAuthenticate *string
}

type UnauthorizedApplicationProblemPlusJSONResponse struct {
	Body    Problem
	Headers UnauthorizedResponseHeaders
}

type StartChargingSessionRequestObject struct {
	Params StartChargingSessionParams
	Body   *StartChargingSessionJSONRequestBody
}

type StartChargingSessionResponseObject interface {
	VisitStartChargingSessionResponse(w http.ResponseWriter) error
}

type StartChargingSession201ResponseHeaders struct {
	Location *string
}

type StartChargingSession201JSONResponse struct {
	Body    ChargingSession
	Headers StartChargingSession201ResponseHeaders
}

func (response StartChargingSession201JSONResponse) VisitStartChargingSessionResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.Body); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	if response.Headers.Location != nil {
		{
			v, err := runtime.StyleParamWithOptions("simple", false, "Location", *response.Headers.Location, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "string", Format: "uri-reference"})
			if err != nil {
				return fmt.Errorf("error styling response header Location: %w", err)
			}

			w.Header().Set("Location", v)
		}
	}

	w.WriteHeader(201)
	_, err := buf.WriteTo(w)

	return err
}

type StartChargingSession400ApplicationProblemPlusJSONResponse struct {
	BadRequestApplicationProblemPlusJSONResponse
}

func (response StartChargingSession400ApplicationProblemPlusJSONResponse) VisitStartChargingSessionResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.BadRequestApplicationProblemPlusJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(400)
	_, err := buf.WriteTo(w)

	return err
}

type StartChargingSession401ApplicationProblemPlusJSONResponse struct {
	UnauthorizedApplicationProblemPlusJSONResponse
}

func (response StartChargingSession401ApplicationProblemPlusJSONResponse) VisitStartChargingSessionResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.Body); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")

	if response.Headers.WWWAuthenticate != nil {
		{
			v, err := runtime.StyleParamWithOptions("simple", false, "WWW-Authenticate", *response.Headers.WWWAuthenticate, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "string", Format: ""})
			if err != nil {
				return fmt.Errorf("error styling response header WWW-Authenticate: %w", err)
			}

			w.Header().Set("WWW-Authenticate", v)
		}
	}

	w.WriteHeader(401)
	_, err := buf.WriteTo(w)

	return err
}

type StartChargingSession403ApplicationProblemPlusJSONResponse struct {
	ForbiddenApplicationProblemPlusJSONResponse
}

func (response StartChargingSession403ApplicationProblemPlusJSONResponse) VisitStartChargingSessionResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.ForbiddenApplicationProblemPlusJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(403)
	_, err := buf.WriteTo(w)

	return err
}

type StartChargingSession409ApplicationProblemPlusJSONResponse Problem

func (response StartChargingSession409ApplicationProblemPlusJSONResponse) VisitStartChargingSessionResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(409)
	_, err := buf.WriteTo(w)

	return err
}

type GetChargingSessionRequestObject struct {
	SessionID SessionID `json:"sessionId"`
}

type GetChargingSessionResponseObject interface {
	VisitGetChargingSessionResponse(w http.ResponseWriter) error
}

type GetChargingSession200JSONResponse ChargingSession

func (response GetChargingSession200JSONResponse) VisitGetChargingSessionResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err := buf.WriteTo(w)

	return err
}

type GetChargingSession401ApplicationProblemPlusJSONResponse struct {
	UnauthorizedApplicationProblemPlusJSONResponse
}

func (response GetChargingSession401ApplicationProblemPlusJSONResponse) VisitGetChargingSessionResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.Body); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")

	if response.Headers.WWWAuthenticate != nil {
		{
			v, err := runtime.StyleParamWithOptions("simple", false, "WWW-Authenticate", *response.Headers.WWWAuthenticate, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "string", Format: ""})
			if err != nil {
				return fmt.Errorf("error styling response header WWW-Authenticate: %w", err)
			}

			w.Header().Set("WWW-Authenticate", v)
		}
	}

	w.WriteHeader(401)
	_, err := buf.WriteTo(w)

	return err
}

type GetChargingSession403ApplicationProblemPlusJSONResponse struct {
	ForbiddenApplicationProblemPlusJSONResponse
}

func (response GetChargingSession403ApplicationProblemPlusJSONResponse) VisitGetChargingSessionResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.ForbiddenApplicationProblemPlusJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(403)
	_, err := buf.WriteTo(w)

	return err
}

type GetChargingSession404ApplicationProblemPlusJSONResponse struct {
	NotFoundApplicationProblemPlusJSONResponse
}

func (response GetChargingSession404ApplicationProblemPlusJSONResponse) VisitGetChargingSessionResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.NotFoundApplicationProblemPlusJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(404)
	_, err := buf.WriteTo(w)

	return err
}

type StopChargingSessionRequestObject struct {
	SessionID SessionID `json:"sessionId"`
	Params    StopChargingSessionParams
}

type StopChargingSessionResponseObject interface {
	VisitStopChargingSessionResponse(w http.ResponseWriter) error
}

type StopChargingSession202JSONResponse ChargingSession

func (response StopChargingSession202JSONResponse) VisitStopChargingSessionResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(202)
	_, err := buf.WriteTo(w)

	return err
}

type StopChargingSession401ApplicationProblemPlusJSONResponse struct {
	UnauthorizedApplicationProblemPlusJSONResponse
}

func (response StopChargingSession401ApplicationProblemPlusJSONResponse) VisitStopChargingSessionResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.Body); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")

	if response.Headers.WWWAuthenticate != nil {
		{
			v, err := runtime.StyleParamWithOptions("simple", false, "WWW-Authenticate", *response.Headers.WWWAuthenticate, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "string", Format: ""})
			if err != nil {
				return fmt.Errorf("error styling response header WWW-Authenticate: %w", err)
			}

			w.Header().Set("WWW-Authenticate", v)
		}
	}

	w.WriteHeader(401)
	_, err := buf.WriteTo(w)

	return err
}

type StopChargingSession403ApplicationProblemPlusJSONResponse struct {
	ForbiddenApplicationProblemPlusJSONResponse
}

func (response StopChargingSession403ApplicationProblemPlusJSONResponse) VisitStopChargingSessionResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.ForbiddenApplicationProblemPlusJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(403)
	_, err := buf.WriteTo(w)

	return err
}

type StopChargingSession404ApplicationProblemPlusJSONResponse struct {
	NotFoundApplicationProblemPlusJSONResponse
}

func (response StopChargingSession404ApplicationProblemPlusJSONResponse) VisitStopChargingSessionResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.NotFoundApplicationProblemPlusJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(404)
	_, err := buf.WriteTo(w)

	return err
}

type GetCurrentDriverRequestObject struct{}

type GetCurrentDriverResponseObject interface {
	VisitGetCurrentDriverResponse(w http.ResponseWriter) error
}

type GetCurrentDriver200JSONResponse Driver

func (response GetCurrentDriver200JSONResponse) VisitGetCurrentDriverResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err := buf.WriteTo(w)

	return err
}

type GetCurrentDriver401ApplicationProblemPlusJSONResponse struct {
	UnauthorizedApplicationProblemPlusJSONResponse
}

func (response GetCurrentDriver401ApplicationProblemPlusJSONResponse) VisitGetCurrentDriverResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.Body); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")

	if response.Headers.WWWAuthenticate != nil {
		{
			v, err := runtime.StyleParamWithOptions("simple", false, "WWW-Authenticate", *response.Headers.WWWAuthenticate, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "string", Format: ""})
			if err != nil {
				return fmt.Errorf("error styling response header WWW-Authenticate: %w", err)
			}

			w.Header().Set("WWW-Authenticate", v)
		}
	}

	w.WriteHeader(401)
	_, err := buf.WriteTo(w)

	return err
}

type GetCurrentDriver403ApplicationProblemPlusJSONResponse struct {
	ForbiddenApplicationProblemPlusJSONResponse
}

func (response GetCurrentDriver403ApplicationProblemPlusJSONResponse) VisitGetCurrentDriverResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.ForbiddenApplicationProblemPlusJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(403)
	_, err := buf.WriteTo(w)

	return err
}

type ListFleetInvoicesRequestObject struct {
	FleetID FleetID `json:"fleetId"`
	Params  ListFleetInvoicesParams
}

type ListFleetInvoicesResponseObject interface {
	VisitListFleetInvoicesResponse(w http.ResponseWriter) error
}

type ListFleetInvoices200JSONResponse []Invoice

func (response ListFleetInvoices200JSONResponse) VisitListFleetInvoicesResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err := buf.WriteTo(w)

	return err
}

type ListFleetInvoices401ApplicationProblemPlusJSONResponse struct {
	UnauthorizedApplicationProblemPlusJSONResponse
}

func (response ListFleetInvoices401ApplicationProblemPlusJSONResponse) VisitListFleetInvoicesResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.Body); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")

	if response.Headers.WWWAuthenticate != nil {
		{
			v, err := runtime.StyleParamWithOptions("simple", false, "WWW-Authenticate", *response.Headers.WWWAuthenticate, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "string", Format: ""})
			if err != nil {
				return fmt.Errorf("error styling response header WWW-Authenticate: %w", err)
			}

			w.Header().Set("WWW-Authenticate", v)
		}
	}

	w.WriteHeader(401)
	_, err := buf.WriteTo(w)

	return err
}

type ListFleetInvoices403ApplicationProblemPlusJSONResponse struct {
	ForbiddenApplicationProblemPlusJSONResponse
}

func (response ListFleetInvoices403ApplicationProblemPlusJSONResponse) VisitListFleetInvoicesResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.ForbiddenApplicationProblemPlusJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(403)
	_, err := buf.WriteTo(w)

	return err
}

type ListFleetVehiclesRequestObject struct {
	FleetID FleetID `json:"fleetId"`
	Params  ListFleetVehiclesParams
}

type ListFleetVehiclesResponseObject interface {
	VisitListFleetVehiclesResponse(w http.ResponseWriter) error
}

type ListFleetVehicles200JSONResponse struct {
	Items      []Vehicle `json:"items"`
	NextCursor *string   `json:"nextCursor"`
}

func (response ListFleetVehicles200JSONResponse) VisitListFleetVehiclesResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err := buf.WriteTo(w)

	return err
}

type ListFleetVehicles401ApplicationProblemPlusJSONResponse struct {
	UnauthorizedApplicationProblemPlusJSONResponse
}

func (response ListFleetVehicles401ApplicationProblemPlusJSONResponse) VisitListFleetVehiclesResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.Body); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")

	if response.Headers.WWWAuthenticate != nil {
		{
			v, err := runtime.StyleParamWithOptions("simple", false, "WWW-Authenticate", *response.Headers.WWWAuthenticate, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "string", Format: ""})
			if err != nil {
				return fmt.Errorf("error styling response header WWW-Authenticate: %w", err)
			}

			w.Header().Set("WWW-Authenticate", v)
		}
	}

	w.WriteHeader(401)
	_, err := buf.WriteTo(w)

	return err
}

type ListFleetVehicles403ApplicationProblemPlusJSONResponse struct {
	ForbiddenApplicationProblemPlusJSONResponse
}

func (response ListFleetVehicles403ApplicationProblemPlusJSONResponse) VisitListFleetVehiclesResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.ForbiddenApplicationProblemPlusJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(403)
	_, err := buf.WriteTo(w)

	return err
}

type ListFleetVehicles404ApplicationProblemPlusJSONResponse struct {
	NotFoundApplicationProblemPlusJSONResponse
}

func (response ListFleetVehicles404ApplicationProblemPlusJSONResponse) VisitListFleetVehiclesResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.NotFoundApplicationProblemPlusJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(404)
	_, err := buf.WriteTo(w)

	return err
}

type GetHealthRequestObject struct{}

type GetHealthResponseObject interface {
	VisitGetHealthResponse(w http.ResponseWriter) error
}

type GetHealth200JSONResponse struct {
	Status GetHealth200JSONResponseBodyStatus `json:"status"`
}

func (response GetHealth200JSONResponse) VisitGetHealthResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err := buf.WriteTo(w)

	return err
}

type SearchChargingStationsRequestObject struct {
	Params SearchChargingStationsParams
}

type SearchChargingStationsResponseObject interface {
	VisitSearchChargingStationsResponse(w http.ResponseWriter) error
}

type SearchChargingStations200JSONResponse struct {
	Items []Station `json:"items"`
	Total int       `json:"total"`
}

func (response SearchChargingStations200JSONResponse) VisitSearchChargingStationsResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err := buf.WriteTo(w)

	return err
}

type SearchChargingStations401ApplicationProblemPlusJSONResponse struct {
	UnauthorizedApplicationProblemPlusJSONResponse
}

func (response SearchChargingStations401ApplicationProblemPlusJSONResponse) VisitSearchChargingStationsResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.Body); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")

	if response.Headers.WWWAuthenticate != nil {
		{
			v, err := runtime.StyleParamWithOptions("simple", false, "WWW-Authenticate", *response.Headers.WWWAuthenticate, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "string", Format: ""})
			if err != nil {
				return fmt.Errorf("error styling response header WWW-Authenticate: %w", err)
			}

			w.Header().Set("WWW-Authenticate", v)
		}
	}

	w.WriteHeader(401)
	_, err := buf.WriteTo(w)

	return err
}

type SearchChargingStations403ApplicationProblemPlusJSONResponse struct {
	ForbiddenApplicationProblemPlusJSONResponse
}

func (response SearchChargingStations403ApplicationProblemPlusJSONResponse) VisitSearchChargingStationsResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.ForbiddenApplicationProblemPlusJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(403)
	_, err := buf.WriteTo(w)

	return err
}

type SearchChargingStations429ApplicationProblemPlusJSONResponse struct {
	RateLimitedApplicationProblemPlusJSONResponse
}

func (response SearchChargingStations429ApplicationProblemPlusJSONResponse) VisitSearchChargingStationsResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.Body); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/problem+json")

	if response.Headers.RetryAfter != nil {
		{
			v, err := runtime.StyleParamWithOptions("simple", false, "Retry-After", *response.Headers.RetryAfter, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "integer", Format: ""})
			if err != nil {
				return fmt.Errorf("error styling response header Retry-After: %w", err)
			}

			w.Header().Set("Retry-After", v)
		}
	}

	w.WriteHeader(429)
	_, err := buf.WriteTo(w)

	return err
}

type StrictServerInterface interface {
	StartChargingSession(ctx context.Context, request StartChargingSessionRequestObject) (StartChargingSessionResponseObject, error)
	GetChargingSession(ctx context.Context, request GetChargingSessionRequestObject) (GetChargingSessionResponseObject, error)
	StopChargingSession(ctx context.Context, request StopChargingSessionRequestObject) (StopChargingSessionResponseObject, error)
	GetCurrentDriver(ctx context.Context, request GetCurrentDriverRequestObject) (GetCurrentDriverResponseObject, error)
	ListFleetInvoices(ctx context.Context, request ListFleetInvoicesRequestObject) (ListFleetInvoicesResponseObject, error)
	ListFleetVehicles(ctx context.Context, request ListFleetVehiclesRequestObject) (ListFleetVehiclesResponseObject, error)
	GetHealth(ctx context.Context, request GetHealthRequestObject) (GetHealthResponseObject, error)
	SearchChargingStations(ctx context.Context, request SearchChargingStationsRequestObject) (SearchChargingStationsResponseObject, error)
}

type StrictHandlerFunc = func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error)

type StrictMiddlewareFunc = func(f StrictHandlerFunc, operationID string) StrictHandlerFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
	SchemaRegistry           SchemaRegistry
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) (ServerInterface, error) {
	registry, err := NewDefaultSchemaRegistry()
	if err != nil {
		return nil, err
	}

	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
		SchemaRegistry: registry,
	}}, nil
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) (ServerInterface, error) {
	if options.RequestErrorHandlerFunc == nil {
		options.RequestErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	if options.ResponseErrorHandlerFunc == nil {
		options.ResponseErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	if options.SchemaRegistry == nil {
		registry, err := NewDefaultSchemaRegistry()
		if err != nil {
			return nil, err
		}

		options.SchemaRegistry = registry
	}

	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}, nil
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

func (sh *strictHandler) StartChargingSession(w http.ResponseWriter, r *http.Request, params StartChargingSessionParams) {
	var request StartChargingSessionRequestObject
	request.Params = params

	if err := sh.validateJSONBody("StartChargingSession", r); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, err)
		return
	}

	var body StartChargingSessionJSONRequestBody
	if err := runtime.BindJSONBody(r, &body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}

	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.StartChargingSession(ctx, request.(StartChargingSessionRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "StartChargingSession")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(StartChargingSessionResponseObject); ok {
		if err := validResponse.VisitStartChargingSessionResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) GetChargingSession(w http.ResponseWriter, r *http.Request, sessionID SessionID) {
	var request GetChargingSessionRequestObject
	request.SessionID = sessionID

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.GetChargingSession(ctx, request.(GetChargingSessionRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetChargingSession")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetChargingSessionResponseObject); ok {
		if err := validResponse.VisitGetChargingSessionResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) StopChargingSession(w http.ResponseWriter, r *http.Request, sessionID SessionID, params StopChargingSessionParams) {
	var request StopChargingSessionRequestObject
	request.SessionID = sessionID
	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.StopChargingSession(ctx, request.(StopChargingSessionRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "StopChargingSession")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(StopChargingSessionResponseObject); ok {
		if err := validResponse.VisitStopChargingSessionResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) GetCurrentDriver(w http.ResponseWriter, r *http.Request) {
	var request GetCurrentDriverRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.GetCurrentDriver(ctx, request.(GetCurrentDriverRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetCurrentDriver")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetCurrentDriverResponseObject); ok {
		if err := validResponse.VisitGetCurrentDriverResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) ListFleetInvoices(w http.ResponseWriter, r *http.Request, fleetID FleetID, params ListFleetInvoicesParams) {
	var request ListFleetInvoicesRequestObject
	request.FleetID = fleetID
	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.ListFleetInvoices(ctx, request.(ListFleetInvoicesRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ListFleetInvoices")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ListFleetInvoicesResponseObject); ok {
		if err := validResponse.VisitListFleetInvoicesResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) ListFleetVehicles(w http.ResponseWriter, r *http.Request, fleetID FleetID, params ListFleetVehiclesParams) {
	var request ListFleetVehiclesRequestObject
	request.FleetID = fleetID
	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.ListFleetVehicles(ctx, request.(ListFleetVehiclesRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ListFleetVehicles")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ListFleetVehiclesResponseObject); ok {
		if err := validResponse.VisitListFleetVehiclesResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	var request GetHealthRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.GetHealth(ctx, request.(GetHealthRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetHealth")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetHealthResponseObject); ok {
		if err := validResponse.VisitGetHealthResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) SearchChargingStations(w http.ResponseWriter, r *http.Request, params SearchChargingStationsParams) {
	var request SearchChargingStationsRequestObject
	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.SearchChargingStations(ctx, request.(SearchChargingStationsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "SearchChargingStations")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(SearchChargingStationsResponseObject); ok {
		if err := validResponse.VisitSearchChargingStationsResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) validateJSONBody(operationID string, r *http.Request) error {
	registry := sh.options.SchemaRegistry
	if registry == nil {
		return nil
	}

	validator, ok := registry.Validator(operationID)
	if !ok {
		return nil
	}

	if r == nil || r.Body == nil {
		return validator.Validate(bytes.NewReader(nil))
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(body))

	return validator.Validate(bytes.NewReader(body))
}

const StartChargingSessionJSONRequestBodySchema = `{
  "additionalProperties": false,
  "dependentRequired": {
    "costCenter": [
      "fleetId"
    ]
  },
  "properties": {
    "connectorId": {
      "type": "string"
    },
    "costCenter": {
      "maxLength": 64,
      "type": "string"
    },
    "fleetId": {
      "type": "string"
    },
    "stationId": {
      "type": "string"
    },
    "vehicleId": {
      "type": "string"
    }
  },
  "required": [
    "stationId",
    "connectorId",
    "vehicleId"
  ],
  "type": "object"
}`

type SchemaRegistry interface {
	Register(name, source string) error
	Validator(name string) (SchemaValidator, bool)
}

type SchemaValidator interface {
	Validate(data io.Reader) error
}

type InMemorySchemaRegistry struct {
	compiler *jsonschema.Compiler
	schemas  map[string]*jsonschema.Schema
}

func NewInMemorySchemaRegistry() *InMemorySchemaRegistry {
	compiler := jsonschema.NewCompiler()
	compiler.AssertFormat()

	return &InMemorySchemaRegistry{
		compiler: compiler,
		schemas:  make(map[string]*jsonschema.Schema),
	}
}

func NewDefaultSchemaRegistry() (*InMemorySchemaRegistry, error) {
	registry := NewInMemorySchemaRegistry()

	schemas := map[string]string{
		"StartChargingSession": StartChargingSessionJSONRequestBodySchema,
	}
	for name, schema := range schemas {
		if err := registry.Register(name, schema); err != nil {
			return nil, err
		}
	}

	return registry, nil
}

func (r *InMemorySchemaRegistry) Register(operationID, source string) error {
	doc, err := jsonschema.UnmarshalJSON(strings.NewReader(source))
	if err != nil {
		return fmt.Errorf("parse JSON schema %q: invalid JSON document: %w", operationID, err)
	}

	if err := r.compiler.AddResource(operationID, doc); err != nil {
		return fmt.Errorf("register JSON schema %q: unable to add resource: %w", operationID, err)
	}

	schema, err := r.compiler.Compile(operationID)
	if err != nil {
		return fmt.Errorf("compile JSON schema %q: invalid definitions or unresolved references: %w", operationID, err)
	}

	r.schemas[operationID] = schema

	return nil
}

func (r *InMemorySchemaRegistry) Validator(operationID string) (SchemaValidator, bool) {
	schema, exists := r.schemas[operationID]
	if !exists {
		return nil, false
	}

	return &CompiledSchemaValidator{schema: schema}, true
}

type CompiledSchemaValidator struct {
	schema *jsonschema.Schema
}

func (v *CompiledSchemaValidator) Validate(data io.Reader) error {
	value, err := jsonschema.UnmarshalJSON(data)
	if err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if err := v.schema.Validate(value); err != nil {
		validationErr, ok := errors.AsType[*jsonschema.ValidationError](err)
		if !ok {
			return fmt.Errorf("invalid JSON: %w", err)
		}

		message := v.formatErrorMessage(validationErr)

		return errors.New(message)
	}

	return nil
}

func (v *CompiledSchemaValidator) formatErrorMessage(err *jsonschema.ValidationError) (result string) {
	defer func() {
		result = strings.TrimSpace(result)
	}()

	if len(err.Causes) > 0 {
		messages := make([]string, 0, len(err.Causes))
		for _, cause := range err.Causes {
			message := v.formatErrorMessage(cause)
			messages = append(messages, message)
		}

		return strings.Join(messages, "; ")
	}

	field := strings.Join(err.InstanceLocation, ".")
	switch rule := err.ErrorKind.(type) {
	case *kind.Required:
		if len(rule.Missing) == 1 {
			return fmt.Sprintf("%s %s is required", field, rule.Missing[0])
		}

		return fmt.Sprintf("%s missing required fields %s", field, strings.Join(rule.Missing, ", "))
	case *kind.Type:
		if len(rule.Want) == 1 {
			return fmt.Sprintf("%s must be %s (got %s)", field, rule.Want[0], rule.Got)
		}

		return fmt.Sprintf("%s must be one of %s (got %s)", field, strings.Join(rule.Want, ", "), rule.Got)
	case *kind.MinLength:
		if rule.Want == 1 {
			return fmt.Sprintf("%s cannot be empty", field)
		}

		return fmt.Sprintf("%s must contain at least %d characters", field, rule.Want)
	case *kind.MaxLength:
		if rule.Want == 1 {
			return fmt.Sprintf("%s must be a single character", field)
		}

		return fmt.Sprintf("%s must contain no more than %d characters", field, rule.Want)
	case *kind.Minimum:
		return fmt.Sprintf("%s must be at least %s", field, rule.Want.FloatString(2))
	case *kind.Maximum:
		return fmt.Sprintf("%s must not exceed %s", field, rule.Want.FloatString(2))
	case *kind.Pattern:
		return fmt.Sprintf("%s does not match the required pattern", field)
	case *kind.Format:
		return fmt.Sprintf("%s must be a valid %s", field, rule.Want)
	case *kind.Enum:
		values := make([]string, 0, len(rule.Want))
		for _, value := range rule.Want {
			values = append(values, fmt.Sprint(value))
		}

		if len(values) == 1 {
			return fmt.Sprintf("%s must be %s", field, values[0])
		}

		return fmt.Sprintf("%s must be one of %s", field, strings.Join(values, ", "))
	default:
		if err.ErrorKind != nil {
			printer := message.NewPrinter(language.English)
			return fmt.Sprintf("%s %s", field, err.ErrorKind.LocalizedString(printer))
		}

		return fmt.Sprintf("%s is invalid", field)
	}
}
