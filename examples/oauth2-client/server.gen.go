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
	jsonschema "github.com/pb33f/jsonschema/v6"
	"github.com/pb33f/jsonschema/v6/kind"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type ServerInterface interface {
	GetDevices(w http.ResponseWriter, r *http.Request)
	GetLegacyImplicit(w http.ResponseWriter, r *http.Request)
	GetLegacyPassword(w http.ResponseWriter, r *http.Request)
	GetOptionalAuth(w http.ResponseWriter, r *http.Request)
	GetPublic(w http.ResponseWriter, r *http.Request)
	GetReports(w http.ResponseWriter, r *http.Request)
	GetServices(w http.ResponseWriter, r *http.Request)
	ListUsers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
}

type Unimplemented struct{}

func (_ Unimplemented) GetDevices(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) GetLegacyImplicit(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) GetLegacyPassword(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) GetOptionalAuth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) GetPublic(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) GetReports(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) GetServices(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) ListUsers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

func (siw *ServerInterfaceWrapper) GetDevices(w http.ResponseWriter, r *http.Request) {
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetDevices(w, r)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetLegacyImplicit(w http.ResponseWriter, r *http.Request) {
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetLegacyImplicit(w, r)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetLegacyPassword(w http.ResponseWriter, r *http.Request) {
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetLegacyPassword(w, r)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetOptionalAuth(w http.ResponseWriter, r *http.Request) {
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetOptionalAuth(w, r)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetPublic(w http.ResponseWriter, r *http.Request) {
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetPublic(w, r)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetReports(w http.ResponseWriter, r *http.Request) {
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetReports(w, r)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetServices(w http.ResponseWriter, r *http.Request) {
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetServices(w, r)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) ListUsers(w http.ResponseWriter, r *http.Request) {
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListUsers(w, r)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) CreateUser(w http.ResponseWriter, r *http.Request) {
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateUser(w, r)
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
		r.Get(options.BaseURL+"/public", wrapper.GetPublic)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/users", wrapper.ListUsers)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/users", wrapper.CreateUser)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/services", wrapper.GetServices)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/devices", wrapper.GetDevices)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/legacy/password", wrapper.GetLegacyPassword)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/legacy/implicit", wrapper.GetLegacyImplicit)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/reports", wrapper.GetReports)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/optional-auth", wrapper.GetOptionalAuth)
	})

	return r
}

type GetDevicesRequestObject struct{}

type GetDevicesResponseObject interface {
	VisitGetDevicesResponse(w http.ResponseWriter) error
}

type GetDevices200Response struct{}

func (response GetDevices200Response) VisitGetDevicesResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type GetLegacyImplicitRequestObject struct{}

type GetLegacyImplicitResponseObject interface {
	VisitGetLegacyImplicitResponse(w http.ResponseWriter) error
}

type GetLegacyImplicit200Response struct{}

func (response GetLegacyImplicit200Response) VisitGetLegacyImplicitResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type GetLegacyPasswordRequestObject struct{}

type GetLegacyPasswordResponseObject interface {
	VisitGetLegacyPasswordResponse(w http.ResponseWriter) error
}

type GetLegacyPassword200Response struct{}

func (response GetLegacyPassword200Response) VisitGetLegacyPasswordResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type GetOptionalAuthRequestObject struct{}

type GetOptionalAuthResponseObject interface {
	VisitGetOptionalAuthResponse(w http.ResponseWriter) error
}

type GetOptionalAuth200Response struct{}

func (response GetOptionalAuth200Response) VisitGetOptionalAuthResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type GetPublicRequestObject struct{}

type GetPublicResponseObject interface {
	VisitGetPublicResponse(w http.ResponseWriter) error
}

type GetPublic200Response struct{}

func (response GetPublic200Response) VisitGetPublicResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type GetReportsRequestObject struct{}

type GetReportsResponseObject interface {
	VisitGetReportsResponse(w http.ResponseWriter) error
}

type GetReports200Response struct{}

func (response GetReports200Response) VisitGetReportsResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type GetServicesRequestObject struct{}

type GetServicesResponseObject interface {
	VisitGetServicesResponse(w http.ResponseWriter) error
}

type GetServices200Response struct{}

func (response GetServices200Response) VisitGetServicesResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type ListUsersRequestObject struct{}

type ListUsersResponseObject interface {
	VisitListUsersResponse(w http.ResponseWriter) error
}

type ListUsers200Response struct{}

func (response ListUsers200Response) VisitListUsersResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type ListUsers401Response struct{}

func (response ListUsers401Response) VisitListUsersResponse(w http.ResponseWriter) error {
	w.WriteHeader(401)
	return nil
}

type ListUsers403Response struct{}

func (response ListUsers403Response) VisitListUsersResponse(w http.ResponseWriter) error {
	w.WriteHeader(403)
	return nil
}

type CreateUserRequestObject struct{}

type CreateUserResponseObject interface {
	VisitCreateUserResponse(w http.ResponseWriter) error
}

type CreateUser201Response struct{}

func (response CreateUser201Response) VisitCreateUserResponse(w http.ResponseWriter) error {
	w.WriteHeader(201)
	return nil
}

type CreateUser401Response struct{}

func (response CreateUser401Response) VisitCreateUserResponse(w http.ResponseWriter) error {
	w.WriteHeader(401)
	return nil
}

type StrictServerInterface interface {
	GetDevices(ctx context.Context, request GetDevicesRequestObject) (GetDevicesResponseObject, error)
	GetLegacyImplicit(ctx context.Context, request GetLegacyImplicitRequestObject) (GetLegacyImplicitResponseObject, error)
	GetLegacyPassword(ctx context.Context, request GetLegacyPasswordRequestObject) (GetLegacyPasswordResponseObject, error)
	GetOptionalAuth(ctx context.Context, request GetOptionalAuthRequestObject) (GetOptionalAuthResponseObject, error)
	GetPublic(ctx context.Context, request GetPublicRequestObject) (GetPublicResponseObject, error)
	GetReports(ctx context.Context, request GetReportsRequestObject) (GetReportsResponseObject, error)
	GetServices(ctx context.Context, request GetServicesRequestObject) (GetServicesResponseObject, error)
	ListUsers(ctx context.Context, request ListUsersRequestObject) (ListUsersResponseObject, error)
	CreateUser(ctx context.Context, request CreateUserRequestObject) (CreateUserResponseObject, error)
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

func (sh *strictHandler) GetDevices(w http.ResponseWriter, r *http.Request) {
	var request GetDevicesRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.GetDevices(ctx, request.(GetDevicesRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetDevices")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetDevicesResponseObject); ok {
		if err := validResponse.VisitGetDevicesResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) GetLegacyImplicit(w http.ResponseWriter, r *http.Request) {
	var request GetLegacyImplicitRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.GetLegacyImplicit(ctx, request.(GetLegacyImplicitRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetLegacyImplicit")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetLegacyImplicitResponseObject); ok {
		if err := validResponse.VisitGetLegacyImplicitResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) GetLegacyPassword(w http.ResponseWriter, r *http.Request) {
	var request GetLegacyPasswordRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.GetLegacyPassword(ctx, request.(GetLegacyPasswordRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetLegacyPassword")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetLegacyPasswordResponseObject); ok {
		if err := validResponse.VisitGetLegacyPasswordResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) GetOptionalAuth(w http.ResponseWriter, r *http.Request) {
	var request GetOptionalAuthRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.GetOptionalAuth(ctx, request.(GetOptionalAuthRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetOptionalAuth")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetOptionalAuthResponseObject); ok {
		if err := validResponse.VisitGetOptionalAuthResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) GetPublic(w http.ResponseWriter, r *http.Request) {
	var request GetPublicRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.GetPublic(ctx, request.(GetPublicRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetPublic")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetPublicResponseObject); ok {
		if err := validResponse.VisitGetPublicResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) GetReports(w http.ResponseWriter, r *http.Request) {
	var request GetReportsRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.GetReports(ctx, request.(GetReportsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetReports")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetReportsResponseObject); ok {
		if err := validResponse.VisitGetReportsResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) GetServices(w http.ResponseWriter, r *http.Request) {
	var request GetServicesRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.GetServices(ctx, request.(GetServicesRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetServices")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetServicesResponseObject); ok {
		if err := validResponse.VisitGetServicesResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	var request ListUsersRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.ListUsers(ctx, request.(ListUsersRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ListUsers")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ListUsersResponseObject); ok {
		if err := validResponse.VisitListUsersResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request CreateUserRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.CreateUser(ctx, request.(CreateUserRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateUser")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CreateUserResponseObject); ok {
		if err := validResponse.VisitCreateUserResponse(w); err != nil {
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

	schemas := map[string]string{}
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
