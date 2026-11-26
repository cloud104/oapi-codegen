package examples

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	chi "github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	jsonschema "github.com/pb33f/jsonschema/v6"
	"github.com/pb33f/jsonschema/v6/kind"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type ServerInterface interface {
	ListBeers(w http.ResponseWriter, r *http.Request, params ListBeersParams)
	CreateBeer(w http.ResponseWriter, r *http.Request)
	DeleteBeer(w http.ResponseWriter, r *http.Request, beerID int)
	GetBeer(w http.ResponseWriter, r *http.Request, beerID int)
	UpdateBeer(w http.ResponseWriter, r *http.Request, beerID int)
	ReplaceBeer(w http.ResponseWriter, r *http.Request, beerID int)
}

type Unimplemented struct{}

func (_ Unimplemented) ListBeers(w http.ResponseWriter, r *http.Request, params ListBeersParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) CreateBeer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) DeleteBeer(w http.ResponseWriter, r *http.Request, beerID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) GetBeer(w http.ResponseWriter, r *http.Request, beerID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) UpdateBeer(w http.ResponseWriter, r *http.Request, beerID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) ReplaceBeer(w http.ResponseWriter, r *http.Request, beerID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

func (siw *ServerInterfaceWrapper) ListBeers(w http.ResponseWriter, r *http.Request) {
	var err error
	_ = err
	var params ListBeersParams

	err = runtime.BindQueryParameterWithOptions("form", true, false, "page", r.URL.Query(), &params.Page, runtime.BindQueryParameterOptions{Type: "integer", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "page"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "page", Err: err})
		}

		return
	}

	err = runtime.BindQueryParameterWithOptions("form", true, false, "pageSize", r.URL.Query(), &params.PageSize, runtime.BindQueryParameterOptions{Type: "integer", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "pageSize"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "pageSize", Err: err})
		}

		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListBeers(w, r, params)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) CreateBeer(w http.ResponseWriter, r *http.Request) {
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateBeer(w, r)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) DeleteBeer(w http.ResponseWriter, r *http.Request) {
	var err error
	_ = err
	var beerID int

	err = runtime.BindStyledParameterWithOptions("simple", "beerId", chi.URLParam(r, "beerId"), &beerID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "integer", Format: "", ValueIsUnescaped: r.URL.RawPath == ""})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "beerId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteBeer(w, r, beerID)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetBeer(w http.ResponseWriter, r *http.Request) {
	var err error
	_ = err
	var beerID int

	err = runtime.BindStyledParameterWithOptions("simple", "beerId", chi.URLParam(r, "beerId"), &beerID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "integer", Format: "", ValueIsUnescaped: r.URL.RawPath == ""})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "beerId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetBeer(w, r, beerID)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) UpdateBeer(w http.ResponseWriter, r *http.Request) {
	var err error
	_ = err
	var beerID int

	err = runtime.BindStyledParameterWithOptions("simple", "beerId", chi.URLParam(r, "beerId"), &beerID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "integer", Format: "", ValueIsUnescaped: r.URL.RawPath == ""})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "beerId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdateBeer(w, r, beerID)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) ReplaceBeer(w http.ResponseWriter, r *http.Request) {
	var err error
	_ = err
	var beerID int

	err = runtime.BindStyledParameterWithOptions("simple", "beerId", chi.URLParam(r, "beerId"), &beerID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "integer", Format: "", ValueIsUnescaped: r.URL.RawPath == ""})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "beerId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ReplaceBeer(w, r, beerID)
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
		r.Get(options.BaseURL+"/beers", wrapper.ListBeers)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/beers", wrapper.CreateBeer)
	})
	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/beers/{beerId}", wrapper.DeleteBeer)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/beers/{beerId}", wrapper.GetBeer)
	})
	r.Group(func(r chi.Router) {
		r.Patch(options.BaseURL+"/beers/{beerId}", wrapper.UpdateBeer)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/beers/{beerId}", wrapper.ReplaceBeer)
	})

	return r
}

type BadRequestJSONResponse Error

type ConflictJSONResponse Error

type NotFoundJSONResponse Error

type UnsupportedMediaTypeJSONResponse Error

type ListBeersRequestObject struct {
	Params ListBeersParams
}

type ListBeersResponseObject interface {
	VisitListBeersResponse(w http.ResponseWriter) error
}

type ListBeers200ResponseHeaders struct {
	Link        *string
	XTotalCount *int
}

type ListBeers200JSONResponse struct {
	Body    BeerPage
	Headers ListBeers200ResponseHeaders
}

func (response ListBeers200JSONResponse) VisitListBeersResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.Body); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	if response.Headers.Link != nil {
		{
			v, err := runtime.StyleParamWithOptions("simple", false, "Link", *response.Headers.Link, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "string", Format: ""})
			if err != nil {
				return fmt.Errorf("error styling response header Link: %w", err)
			}

			w.Header().Set("Link", v)
		}
	}

	if response.Headers.XTotalCount != nil {
		{
			v, err := runtime.StyleParamWithOptions("simple", false, "X-Total-Count", *response.Headers.XTotalCount, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "integer", Format: ""})
			if err != nil {
				return fmt.Errorf("error styling response header X-Total-Count: %w", err)
			}

			w.Header().Set("X-Total-Count", v)
		}
	}

	w.WriteHeader(200)
	_, err := buf.WriteTo(w)

	return err
}

type ListBeers200ApplicationXMLResponse struct {
	Body          io.Reader
	Headers       ListBeers200ResponseHeaders
	ContentLength int64
}

func (response ListBeers200ApplicationXMLResponse) VisitListBeersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/xml")

	if response.ContentLength != 0 {
		w.Header().Set("Content-Length", fmt.Sprint(response.ContentLength))
	}

	if response.Headers.Link != nil {
		{
			v, err := runtime.StyleParamWithOptions("simple", false, "Link", *response.Headers.Link, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "string", Format: ""})
			if err != nil {
				return fmt.Errorf("error styling response header Link: %w", err)
			}

			w.Header().Set("Link", v)
		}
	}

	if response.Headers.XTotalCount != nil {
		{
			v, err := runtime.StyleParamWithOptions("simple", false, "X-Total-Count", *response.Headers.XTotalCount, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "integer", Format: ""})
			if err != nil {
				return fmt.Errorf("error styling response header X-Total-Count: %w", err)
			}

			w.Header().Set("X-Total-Count", v)
		}
	}

	w.WriteHeader(200)

	if closer, ok := response.Body.(io.ReadCloser); ok {
		defer closer.Close()
	}

	_, err := io.Copy(w, response.Body)

	return err
}

type ListBeers202ResponseHeaders struct {
	Link        *string
	RetryAfter  *int
	XRequestID  *string
	XTotalCount int
}

type ListBeers202ApplicationXMLResponse struct {
	Body          io.Reader
	Headers       ListBeers202ResponseHeaders
	ContentLength int64
}

func (response ListBeers202ApplicationXMLResponse) VisitListBeersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/xml")

	if response.ContentLength != 0 {
		w.Header().Set("Content-Length", fmt.Sprint(response.ContentLength))
	}

	if response.Headers.Link != nil {
		{
			v, err := runtime.StyleParamWithOptions("simple", false, "Link", *response.Headers.Link, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "string", Format: ""})
			if err != nil {
				return fmt.Errorf("error styling response header Link: %w", err)
			}

			w.Header().Set("Link", v)
		}
	}

	if response.Headers.RetryAfter != nil {
		{
			v, err := runtime.StyleParamWithOptions("simple", false, "Retry-After", *response.Headers.RetryAfter, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "integer", Format: ""})
			if err != nil {
				return fmt.Errorf("error styling response header Retry-After: %w", err)
			}

			w.Header().Set("Retry-After", v)
		}
	}

	if response.Headers.XRequestID != nil {
		{
			v, err := runtime.StyleParamWithOptions("simple", false, "X-Request-Id", *response.Headers.XRequestID, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "string", Format: ""})
			if err != nil {
				return fmt.Errorf("error styling response header X-Request-Id: %w", err)
			}

			w.Header().Set("X-Request-Id", v)
		}
	}
	{
		v, err := runtime.StyleParamWithOptions("simple", false, "X-Total-Count", response.Headers.XTotalCount, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "integer", Format: ""})
		if err != nil {
			return fmt.Errorf("error styling response header X-Total-Count: %w", err)
		}

		w.Header().Set("X-Total-Count", v)
	}

	w.WriteHeader(202)

	if closer, ok := response.Body.(io.ReadCloser); ok {
		defer closer.Close()
	}

	_, err := io.Copy(w, response.Body)

	return err
}

type CreateBeerRequestObject struct {
	JSONBody      *CreateBeerJSONRequestBody
	FormdataBody  *CreateBeerFormdataRequestBody
	MultipartBody *multipart.Reader
	TextBody      *CreateBeerTextRequestBody
}

type CreateBeerResponseObject interface {
	VisitCreateBeerResponse(w http.ResponseWriter) error
}

type CreateBeer201JSONResponse Beer

func (response CreateBeer201JSONResponse) VisitCreateBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	_, err := buf.WriteTo(w)

	return err
}

type CreateBeer400JSONResponse struct{ BadRequestJSONResponse }

func (response CreateBeer400JSONResponse) VisitCreateBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.BadRequestJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	_, err := buf.WriteTo(w)

	return err
}

type CreateBeer409JSONResponse struct{ ConflictJSONResponse }

func (response CreateBeer409JSONResponse) VisitCreateBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.ConflictJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)
	_, err := buf.WriteTo(w)

	return err
}

type CreateBeer415JSONResponse struct {
	UnsupportedMediaTypeJSONResponse
}

func (response CreateBeer415JSONResponse) VisitCreateBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.UnsupportedMediaTypeJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(415)
	_, err := buf.WriteTo(w)

	return err
}

type DeleteBeerRequestObject struct {
	BeerID int `json:"beerId"`
}

type DeleteBeerResponseObject interface {
	VisitDeleteBeerResponse(w http.ResponseWriter) error
}

type DeleteBeer204Response struct{}

func (response DeleteBeer204Response) VisitDeleteBeerResponse(w http.ResponseWriter) error {
	w.WriteHeader(204)
	return nil
}

type DeleteBeer404JSONResponse struct{ NotFoundJSONResponse }

func (response DeleteBeer404JSONResponse) VisitDeleteBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.NotFoundJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)
	_, err := buf.WriteTo(w)

	return err
}

type DeleteBeer409JSONResponse struct{ ConflictJSONResponse }

func (response DeleteBeer409JSONResponse) VisitDeleteBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.ConflictJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)
	_, err := buf.WriteTo(w)

	return err
}

type DeleteBeerdefaultJSONResponse struct {
	Body       map[string]any
	StatusCode int
}

func (response DeleteBeerdefaultJSONResponse) VisitDeleteBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.Body); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	_, err := buf.WriteTo(w)

	return err
}

type GetBeerRequestObject struct {
	BeerID int `json:"beerId"`
}

type GetBeerResponseObject interface {
	VisitGetBeerResponse(w http.ResponseWriter) error
}

type GetBeer200JSONResponse struct {
	Abv   *float32 `json:"abv,omitempty"`
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Style *string  `json:"style,omitempty"`
}

func (response GetBeer200JSONResponse) VisitGetBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err := buf.WriteTo(w)

	return err
}

type GetBeer404JSONResponse struct{ NotFoundJSONResponse }

func (response GetBeer404JSONResponse) VisitGetBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.NotFoundJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)
	_, err := buf.WriteTo(w)

	return err
}

type UpdateBeerRequestObject struct {
	BeerID int `json:"beerId"`
	Body   *UpdateBeerJSONRequestBody
}

type UpdateBeerResponseObject interface {
	VisitUpdateBeerResponse(w http.ResponseWriter) error
}

type UpdateBeer200JSONResponse Beer

func (response UpdateBeer200JSONResponse) VisitUpdateBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err := buf.WriteTo(w)

	return err
}

type UpdateBeer400JSONResponse struct{ BadRequestJSONResponse }

func (response UpdateBeer400JSONResponse) VisitUpdateBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.BadRequestJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	_, err := buf.WriteTo(w)

	return err
}

type UpdateBeer404JSONResponse struct{ NotFoundJSONResponse }

func (response UpdateBeer404JSONResponse) VisitUpdateBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.NotFoundJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)
	_, err := buf.WriteTo(w)

	return err
}

type UpdateBeer409JSONResponse struct{ ConflictJSONResponse }

func (response UpdateBeer409JSONResponse) VisitUpdateBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.ConflictJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)
	_, err := buf.WriteTo(w)

	return err
}

type UpdateBeer415JSONResponse struct {
	UnsupportedMediaTypeJSONResponse
}

func (response UpdateBeer415JSONResponse) VisitUpdateBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.UnsupportedMediaTypeJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(415)
	_, err := buf.WriteTo(w)

	return err
}

type ReplaceBeerRequestObject struct {
	BeerID int `json:"beerId"`
	Body   *ReplaceBeerJSONRequestBody
}

type ReplaceBeerResponseObject interface {
	VisitReplaceBeerResponse(w http.ResponseWriter) error
}

type ReplaceBeer200JSONResponse Beer

func (response ReplaceBeer200JSONResponse) VisitReplaceBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err := buf.WriteTo(w)

	return err
}

type ReplaceBeer400JSONResponse struct{ BadRequestJSONResponse }

func (response ReplaceBeer400JSONResponse) VisitReplaceBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.BadRequestJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	_, err := buf.WriteTo(w)

	return err
}

type ReplaceBeer404JSONResponse struct{ NotFoundJSONResponse }

func (response ReplaceBeer404JSONResponse) VisitReplaceBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.NotFoundJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)
	_, err := buf.WriteTo(w)

	return err
}

type ReplaceBeer409JSONResponse struct{ ConflictJSONResponse }

func (response ReplaceBeer409JSONResponse) VisitReplaceBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.ConflictJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)
	_, err := buf.WriteTo(w)

	return err
}

type ReplaceBeer415JSONResponse struct {
	UnsupportedMediaTypeJSONResponse
}

func (response ReplaceBeer415JSONResponse) VisitReplaceBeerResponse(w http.ResponseWriter) error {
	var buf bytes.Buffer
	if err := runtime.JSONMarshalWrite(&buf, response.UnsupportedMediaTypeJSONResponse); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(415)
	_, err := buf.WriteTo(w)

	return err
}

type StrictServerInterface interface {
	ListBeers(ctx context.Context, request ListBeersRequestObject) (ListBeersResponseObject, error)
	CreateBeer(ctx context.Context, request CreateBeerRequestObject) (CreateBeerResponseObject, error)
	DeleteBeer(ctx context.Context, request DeleteBeerRequestObject) (DeleteBeerResponseObject, error)
	GetBeer(ctx context.Context, request GetBeerRequestObject) (GetBeerResponseObject, error)
	UpdateBeer(ctx context.Context, request UpdateBeerRequestObject) (UpdateBeerResponseObject, error)
	ReplaceBeer(ctx context.Context, request ReplaceBeerRequestObject) (ReplaceBeerResponseObject, error)
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

func (sh *strictHandler) ListBeers(w http.ResponseWriter, r *http.Request, params ListBeersParams) {
	var request ListBeersRequestObject
	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.ListBeers(ctx, request.(ListBeersRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ListBeers")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ListBeersResponseObject); ok {
		if err := validResponse.VisitListBeersResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) CreateBeer(w http.ResponseWriter, r *http.Request) {
	var request CreateBeerRequestObject

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		if err := sh.validateJSONBody("CreateBeer", r); err != nil {
			sh.options.RequestErrorHandlerFunc(w, r, err)
			return
		}

		var body CreateBeerJSONRequestBody
		if err := runtime.BindJSONBody(r, &body); err != nil {
			sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
			return
		}

		request.JSONBody = &body
	}

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode formdata: %w", err))
			return
		}

		var body CreateBeerFormdataRequestBody
		if err := runtime.BindForm(&body, r.Form, nil, nil); err != nil {
			sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't bind formdata: %w", err))
			return
		}

		request.FormdataBody = &body
	}

	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		if reader, err := r.MultipartReader(); err != nil {
			sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode multipart body: %w", err))
			return
		} else {
			request.MultipartBody = reader
		}
	}

	if strings.HasPrefix(r.Header.Get("Content-Type"), "text/plain") {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't read body: %w", err))
			return
		}

		body := CreateBeerTextRequestBody(data)
		request.TextBody = &body
	}

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.CreateBeer(ctx, request.(CreateBeerRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateBeer")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CreateBeerResponseObject); ok {
		if err := validResponse.VisitCreateBeerResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) DeleteBeer(w http.ResponseWriter, r *http.Request, beerID int) {
	var request DeleteBeerRequestObject
	request.BeerID = beerID

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.DeleteBeer(ctx, request.(DeleteBeerRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "DeleteBeer")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(DeleteBeerResponseObject); ok {
		if err := validResponse.VisitDeleteBeerResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) GetBeer(w http.ResponseWriter, r *http.Request, beerID int) {
	var request GetBeerRequestObject
	request.BeerID = beerID

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.GetBeer(ctx, request.(GetBeerRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetBeer")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetBeerResponseObject); ok {
		if err := validResponse.VisitGetBeerResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) UpdateBeer(w http.ResponseWriter, r *http.Request, beerID int) {
	var request UpdateBeerRequestObject
	request.BeerID = beerID

	if err := sh.validateJSONBody("UpdateBeer", r); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, err)
		return
	}

	var body UpdateBeerJSONRequestBody
	if err := runtime.BindJSONBody(r, &body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}

	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.UpdateBeer(ctx, request.(UpdateBeerRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "UpdateBeer")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(UpdateBeerResponseObject); ok {
		if err := validResponse.VisitUpdateBeerResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) ReplaceBeer(w http.ResponseWriter, r *http.Request, beerID int) {
	var request ReplaceBeerRequestObject
	request.BeerID = beerID

	if err := sh.validateJSONBody("ReplaceBeer", r); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, err)
		return
	}

	var body ReplaceBeerJSONRequestBody
	if err := runtime.BindJSONBody(r, &body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}

	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return sh.ssi.ReplaceBeer(ctx, request.(ReplaceBeerRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ReplaceBeer")
	}

	response, err := handler(r.Context(), w, r, request)
	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ReplaceBeerResponseObject); ok {
		if err := validResponse.VisitReplaceBeerResponse(w); err != nil {
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

const CreateBeerJSONRequestBodySchema = `{
  "additionalProperties": false,
  "properties": {
    "abv": {
      "format": "float",
      "maximum": 100,
      "minimum": 0,
      "type": "number"
    },
    "brewery": {
      "maxLength": 200,
      "minLength": 1,
      "type": "string"
    },
    "ibu": {
      "maximum": 200,
      "minimum": 0,
      "type": "integer"
    },
    "name": {
      "maxLength": 100,
      "minLength": 1,
      "type": "string"
    },
    "style": {
      "maxLength": 100,
      "minLength": 1,
      "type": "string"
    }
  },
  "required": [
    "name",
    "style",
    "abv"
  ],
  "type": "object"
}`

const UpdateBeerJSONRequestBodySchema = `{
  "additionalProperties": false,
  "minProperties": 1,
  "properties": {
    "abv": {
      "format": "float",
      "maximum": 100,
      "minimum": 0,
      "type": "number"
    },
    "brewery": {
      "maxLength": 200,
      "minLength": 1,
      "type": "string"
    },
    "ibu": {
      "maximum": 200,
      "minimum": 0,
      "type": "integer"
    },
    "name": {
      "maxLength": 100,
      "minLength": 1,
      "type": "string"
    },
    "style": {
      "maxLength": 100,
      "minLength": 1,
      "type": "string"
    }
  },
  "type": "object"
}`

const ReplaceBeerJSONRequestBodySchema = `{
  "additionalProperties": false,
  "properties": {
    "abv": {
      "format": "float",
      "maximum": 100,
      "minimum": 0,
      "type": "number"
    },
    "brewery": {
      "maxLength": 200,
      "minLength": 1,
      "type": "string"
    },
    "ibu": {
      "maximum": 200,
      "minimum": 0,
      "type": "integer"
    },
    "name": {
      "maxLength": 100,
      "minLength": 1,
      "type": "string"
    },
    "style": {
      "maxLength": 100,
      "minLength": 1,
      "type": "string"
    }
  },
  "required": [
    "name",
    "style",
    "abv"
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
		"CreateBeer":  CreateBeerJSONRequestBodySchema,
		"UpdateBeer":  UpdateBeerJSONRequestBodySchema,
		"ReplaceBeer": ReplaceBeerJSONRequestBodySchema,
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
