package examples
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)
type ServerInterface interface {
	GetPing(w http.ResponseWriter, r *http.Request)
}
type Unimplemented struct{}
func (_ Unimplemented) GetPing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}
type MiddlewareFunc func(http.Handler) http.Handler
func (siw *ServerInterfaceWrapper) GetPing(w http.ResponseWriter, r *http.Request) {
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetPing(w, r)
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
		r.Get(options.BaseURL+"/ping", wrapper.GetPing)
	})
	return r
}
type GetPingRequestObject struct {
}
type GetPingResponseObject interface {
	VisitGetPingResponse(w http.ResponseWriter) error
}
type GetPing200JSONResponse Pong
func (response GetPing200JSONResponse) VisitGetPingResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err := buf.WriteTo(w)
	return err
}
type StrictServerInterface interface {
	GetPing(ctx context.Context, request GetPingRequestObject) (GetPingResponseObject, error)
}
type StrictHandlerFunc = func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error)
type StrictMiddlewareFunc = func(f StrictHandlerFunc, operationID string) StrictHandlerFunc
type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc      func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc     func(w http.ResponseWriter, r *http.Request, err error)
	RequestBodyValidatorRegistry RequestBodyValidatorRegistry
}
type RequestBodyValidatorRegistry interface {
	Validator(operationID string) (RequestBodyValidator, bool)
}
type RequestBodyValidator interface {
	Validate(data io.Reader) error
}
func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}
func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
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
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}
type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}
func (sh *strictHandler) GetPing(w http.ResponseWriter, r *http.Request) {
	var request GetPingRequestObject
	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.GetPing(ctx, request.(GetPingRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetPing")
	}
	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetPingResponseObject); ok {
		if err := validResponse.VisitGetPingResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}
func (sh *strictHandler) validateJSONBody(operationID string, r *http.Request) error {
	registry := sh.options.RequestBodyValidatorRegistry
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
type JSONSchemaRegistry struct {
	compiler *jsonschema.Compiler
	schemas  map[string]*jsonschema.Schema
}
func NewJSONSchemaRegistry() *JSONSchemaRegistry {
	compiler := jsonschema.NewCompiler()
	compiler.AssertFormat()
	return &JSONSchemaRegistry{
		compiler: compiler,
		schemas:  make(map[string]*jsonschema.Schema),
	}
}
func NewDefaultJSONSchemaRegistry() (*JSONSchemaRegistry, error) {
	registry := NewJSONSchemaRegistry()
	schemas := map[string]string{}
	for name, schema := range schemas {
		if err := registry.Register(name, schema); err != nil {
			return nil, err
		}
	}
	return registry, nil
}
func (r *JSONSchemaRegistry) Register(operationID, source string) error {
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
func (r *JSONSchemaRegistry) Validator(operationID string) (RequestBodyValidator, bool) {
	schema, exists := r.schemas[operationID]
	if !exists {
		return nil, false
	}
	return &JSONSchemaValidator{schema: schema}, true
}
type JSONSchemaValidator struct {
	schema *jsonschema.Schema
}
func (v *JSONSchemaValidator) Validate(data io.Reader) error {
	value, err := jsonschema.UnmarshalJSON(data)
	if err != nil {
		return &InvalidJSONError{Err: err}
	}
	if err := v.schema.Validate(value); err != nil {
		return &JSONSchemaValidationError{Err: err}
	}
	return nil
}
type InvalidJSONError struct {
	Err error
}
func (e *InvalidJSONError) Error() string {
	return fmt.Sprintf("Invalid JSON: %s", e.Err.Error())
}
func (e *InvalidJSONError) Unwrap() error {
	return e.Err
}
type JSONSchemaValidationError struct {
	Err error
}
func (e *JSONSchemaValidationError) Error() string {
	if validationErr, ok := errors.AsType[*jsonschema.ValidationError](e.Err); ok {
		return e.formatErrorMessage(validationErr)
	}
	return fmt.Sprintf("JSON schema validation failed: %s", e.Err.Error())
}
func (e *JSONSchemaValidationError) Unwrap() error {
	return e.Err
}
func (e *JSONSchemaValidationError) formatErrorMessage(err *jsonschema.ValidationError) (result string) {
	defer func() {
		result = strings.TrimSpace(result)
	}()
	if len(err.Causes) > 0 {
		messages := make([]string, 0, len(err.Causes))
		for _, cause := range err.Causes {
			message := e.formatErrorMessage(cause)
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
