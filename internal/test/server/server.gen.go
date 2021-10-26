// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/discord-gophers/goapi-gen version (devel) DO NOT EDIT.
package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/discord-gophers/goapi-gen/pkg/runtime"
	openapi_types "github.com/discord-gophers/goapi-gen/pkg/types"
	"github.com/go-chi/chi/v5"
)

// EveryTypeOptional defines model for EveryTypeOptional.
type EveryTypeOptional struct {
	ArrayInlineField     *[]int              `json:"array_inline_field,omitempty"`
	ArrayReferencedField *[]SomeObject       `json:"array_referenced_field,omitempty"`
	BoolField            *bool               `json:"bool_field,omitempty"`
	ByteField            *[]byte             `json:"byte_field,omitempty"`
	DateField            *openapi_types.Date `json:"date_field,omitempty"`
	DateTimeField        *time.Time          `json:"date_time_field,omitempty"`
	DoubleField          *float64            `json:"double_field,omitempty"`
	FloatField           *float32            `json:"float_field,omitempty"`
	InlineObjectField    *struct {
		Name   string `json:"name"`
		Number int    `json:"number"`
	} `json:"inline_object_field,omitempty"`
	Int32Field      *int32      `json:"int32_field,omitempty"`
	Int64Field      *int64      `json:"int64_field,omitempty"`
	IntField        *int        `json:"int_field,omitempty"`
	NumberField     *float32    `json:"number_field,omitempty"`
	ReferencedField *SomeObject `json:"referenced_field,omitempty"`
	StringField     *string     `json:"string_field,omitempty"`
}

// EveryTypeRequired defines model for EveryTypeRequired.
type EveryTypeRequired struct {
	ArrayInlineField     []int                `json:"array_inline_field"`
	ArrayReferencedField []SomeObject         `json:"array_referenced_field"`
	BoolField            bool                 `json:"bool_field"`
	ByteField            []byte               `json:"byte_field"`
	DateField            openapi_types.Date   `json:"date_field"`
	DateTimeField        time.Time            `json:"date_time_field"`
	DoubleField          float64              `json:"double_field"`
	EmailField           *openapi_types.Email `json:"email_field,omitempty"`
	FloatField           float32              `json:"float_field"`
	InlineObjectField    struct {
		Name   string `json:"name"`
		Number int    `json:"number"`
	} `json:"inline_object_field"`
	Int32Field      int32      `json:"int32_field"`
	Int64Field      int64      `json:"int64_field"`
	IntField        int        `json:"int_field"`
	NumberField     float32    `json:"number_field"`
	ReferencedField SomeObject `json:"referenced_field"`
	StringField     string     `json:"string_field"`
}

// ReservedKeyword defines model for ReservedKeyword.
type ReservedKeyword struct {
	Channel *string `json:"channel,omitempty"`
}

// Resource defines model for Resource.
type Resource struct {
	Name  string  `json:"name"`
	Value float32 `json:"value"`
}

// SomeObject defines model for some_object.
type SomeObject struct {
	Name string `json:"name"`
}

// Argument defines model for argument.
type Argument string

// ResponseWithReference defines model for ResponseWithReference.
type ResponseWithReference SomeObject

// SimpleResponse defines model for SimpleResponse.
type SimpleResponse struct {
	Name string `json:"name"`
}

// GetWithArgsParams defines parameters for GetWithArgs.
type GetWithArgsParams struct {
	// An optional query argument
	OptionalArgument *int64 `json:"optional_argument,omitempty"`

	// An optional query argument
	RequiredArgument int64 `json:"required_argument"`

	// An optional query argument
	HeaderArgument *int32 `json:"header_argument,omitempty"`
}

// GetWithContentTypeParamsContentType defines parameters for GetWithContentType.
type GetWithContentTypeParamsContentType string

// CreateResourceJSONBody defines parameters for CreateResource.
type CreateResourceJSONBody EveryTypeRequired

// CreateResource2JSONBody defines parameters for CreateResource2.
type CreateResource2JSONBody Resource

// CreateResource2Params defines parameters for CreateResource2.
type CreateResource2Params struct {
	// Some query argument
	InlineQueryArgument *int `json:"inline_query_argument,omitempty"`
}

// UpdateResource3JSONBody defines parameters for UpdateResource3.
type UpdateResource3JSONBody struct {
	ID   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// CreateResourceJSONRequestBody defines body for CreateResource for application/json ContentType.
type CreateResourceJSONRequestBody CreateResourceJSONBody

// CreateResource2JSONRequestBody defines body for CreateResource2 for application/json ContentType.
type CreateResource2JSONRequestBody CreateResource2JSONBody

// UpdateResource3JSONRequestBody defines body for UpdateResource3 for application/json ContentType.
type UpdateResource3JSONRequestBody UpdateResource3JSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// get every type optional
	// (GET /every-type-optional)
	GetEveryTypeOptional(w http.ResponseWriter, r *http.Request)
	// Get resource via simple path
	// (GET /get-simple)
	GetSimple(w http.ResponseWriter, r *http.Request)
	// Getter with referenced parameter and referenced response
	// (GET /get-with-args)
	GetWithArgs(w http.ResponseWriter, r *http.Request, params GetWithArgsParams)
	// Getter with referenced parameter and referenced response
	// (GET /get-with-references/{global_argument}/{argument})
	GetWithReferences(w http.ResponseWriter, r *http.Request, globalArgument int64, argument Argument)

	// (GET /get-with-tagged-middleware)
	GetWithTaggedMiddleware(w http.ResponseWriter, r *http.Request)

	// (POST /get-with-tagged-middleware)
	PostWithTaggedMiddleware(w http.ResponseWriter, r *http.Request)
	// Get an object by ID
	// (GET /get-with-type/{content_type})
	GetWithContentType(w http.ResponseWriter, r *http.Request, contentType GetWithContentTypeParamsContentType)
	// get with reserved keyword
	// (GET /reserved-keyword)
	GetReservedKeyword(w http.ResponseWriter, r *http.Request)
	// Create a resource
	// (POST /resource/{argument})
	CreateResource(w http.ResponseWriter, r *http.Request, argument Argument)
	// Create a resource with inline parameter
	// (POST /resource2/{inline_argument})
	CreateResource2(w http.ResponseWriter, r *http.Request, inlineArgument int, params CreateResource2Params)
	// Update a resource with inline body. The parameter name is a reserved
	// keyword, so make sure that gets prefixed to avoid syntax errors
	// (PUT /resource3/{fallthrough})
	UpdateResource3(w http.ResponseWriter, r *http.Request, pFallthrough int)
	// get response with reference
	// (GET /response-with-reference)
	GetResponseWithReference(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	TaggedMiddlewares  map[string]MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetEveryTypeOptional operation middleware
func (siw *ServerInterfaceWrapper) GetEveryTypeOptional(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetEveryTypeOptional(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler).ServeHTTP
	}

	handler(w, r.WithContext(ctx))
}

// GetSimple operation middleware
func (siw *ServerInterfaceWrapper) GetSimple(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetSimple(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler).ServeHTTP
	}

	handler(w, r.WithContext(ctx))
}

// GetWithArgs operation middleware
func (siw *ServerInterfaceWrapper) GetWithArgs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parameter object where we will unmarshal all parameters from the context
	var params GetWithArgsParams

	// ------------- Optional query parameter "optional_argument" -------------
	if paramValue := r.URL.Query().Get("optional_argument"); paramValue != "" {

	}

	if err := runtime.BindQueryParameter("form", true, false, "optional_argument", r.URL.Query(), &params.OptionalArgument); err != nil {
		err = fmt.Errorf("invalid format for parameter optional_argument: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
		return
	}

	// ------------- Required query parameter "required_argument" -------------
	if paramValue := r.URL.Query().Get("required_argument"); paramValue != "" {

	} else {
		err := fmt.Errorf("query argument required_argument is required, but not found")
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{err})
		return
	}

	if err := runtime.BindQueryParameter("form", true, true, "required_argument", r.URL.Query(), &params.RequiredArgument); err != nil {
		err = fmt.Errorf("invalid format for parameter required_argument: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
		return
	}

	headers := r.Header

	// ------------- Optional header parameter "header_argument" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("header_argument")]; found {
		var HeaderArgument int32
		n := len(valueList)
		if n != 1 {
			err := fmt.Errorf("expected one value for header_argument, got %d", n)
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{err})
			return
		}

		if err := runtime.BindStyledParameterWithLocation("simple", false, "header_argument", runtime.ParamLocationHeader, valueList[0], &HeaderArgument); err != nil {
			err = fmt.Errorf("invalid format for parameter header_argument: %w", err)
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
			return
		}

		params.HeaderArgument = &HeaderArgument

	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetWithArgs(w, r, params)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler).ServeHTTP
	}

	handler(w, r.WithContext(ctx))
}

// GetWithReferences operation middleware
func (siw *ServerInterfaceWrapper) GetWithReferences(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "global_argument" -------------
	var globalArgument int64

	if err := runtime.BindStyledParameter("simple", false, "global_argument", chi.URLParam(r, "global_argument"), &globalArgument); err != nil {
		err = fmt.Errorf("invalid format for parameter global_argument: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
		return
	}

	// ------------- Path parameter "argument" -------------
	var argument Argument

	if err := runtime.BindStyledParameter("simple", false, "argument", chi.URLParam(r, "argument"), &argument); err != nil {
		err = fmt.Errorf("invalid format for parameter argument: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetWithReferences(w, r, globalArgument, argument)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler).ServeHTTP
	}

	handler(w, r.WithContext(ctx))
}

// GetWithTaggedMiddleware operation middleware
func (siw *ServerInterfaceWrapper) GetWithTaggedMiddleware(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetWithTaggedMiddleware(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler).ServeHTTP
	}

	// Operation specific middleware
	if siw.TaggedMiddlewares != nil {
		if middleware, ok := siw.TaggedMiddlewares["pathMiddleware"]; ok {
			handler = middleware(handler).ServeHTTP
		}
	}

	handler(w, r.WithContext(ctx))
}

// PostWithTaggedMiddleware operation middleware
func (siw *ServerInterfaceWrapper) PostWithTaggedMiddleware(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostWithTaggedMiddleware(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler).ServeHTTP
	}

	// Operation specific middleware
	if siw.TaggedMiddlewares != nil {
		if middleware, ok := siw.TaggedMiddlewares["pathMiddleware"]; ok {
			handler = middleware(handler).ServeHTTP
		}
		if middleware, ok := siw.TaggedMiddlewares["operationMiddleware"]; ok {
			handler = middleware(handler).ServeHTTP
		}
	}

	handler(w, r.WithContext(ctx))
}

// GetWithContentType operation middleware
func (siw *ServerInterfaceWrapper) GetWithContentType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "content_type" -------------
	var contentType GetWithContentTypeParamsContentType

	if err := runtime.BindStyledParameter("simple", false, "content_type", chi.URLParam(r, "content_type"), &contentType); err != nil {
		err = fmt.Errorf("invalid format for parameter content_type: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetWithContentType(w, r, contentType)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler).ServeHTTP
	}

	handler(w, r.WithContext(ctx))
}

// GetReservedKeyword operation middleware
func (siw *ServerInterfaceWrapper) GetReservedKeyword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetReservedKeyword(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler).ServeHTTP
	}

	handler(w, r.WithContext(ctx))
}

// CreateResource operation middleware
func (siw *ServerInterfaceWrapper) CreateResource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "argument" -------------
	var argument Argument

	if err := runtime.BindStyledParameter("simple", false, "argument", chi.URLParam(r, "argument"), &argument); err != nil {
		err = fmt.Errorf("invalid format for parameter argument: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateResource(w, r, argument)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler).ServeHTTP
	}

	handler(w, r.WithContext(ctx))
}

// CreateResource2 operation middleware
func (siw *ServerInterfaceWrapper) CreateResource2(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "inline_argument" -------------
	var inlineArgument int

	if err := runtime.BindStyledParameter("simple", false, "inline_argument", chi.URLParam(r, "inline_argument"), &inlineArgument); err != nil {
		err = fmt.Errorf("invalid format for parameter inline_argument: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params CreateResource2Params

	// ------------- Optional query parameter "inline_query_argument" -------------
	if paramValue := r.URL.Query().Get("inline_query_argument"); paramValue != "" {

	}

	if err := runtime.BindQueryParameter("form", true, false, "inline_query_argument", r.URL.Query(), &params.InlineQueryArgument); err != nil {
		err = fmt.Errorf("invalid format for parameter inline_query_argument: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateResource2(w, r, inlineArgument, params)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler).ServeHTTP
	}

	handler(w, r.WithContext(ctx))
}

// UpdateResource3 operation middleware
func (siw *ServerInterfaceWrapper) UpdateResource3(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "fallthrough" -------------
	var pFallthrough int

	if err := runtime.BindStyledParameter("simple", false, "fallthrough", chi.URLParam(r, "fallthrough"), &pFallthrough); err != nil {
		err = fmt.Errorf("invalid format for parameter fallthrough: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdateResource3(w, r, pFallthrough)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler).ServeHTTP
	}

	handler(w, r.WithContext(ctx))
}

// GetResponseWithReference operation middleware
func (siw *ServerInterfaceWrapper) GetResponseWithReference(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetResponseWithReference(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler).ServeHTTP
	}

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	error
}
type UnmarshalingParamError struct {
	error
}
type RequiredParamError struct {
	error
}
type RequiredHeaderError struct {
	error
}
type InvalidParamFormatError struct {
	error
}
type TooManyValuesForParamError struct {
	error
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL           string
	BaseRouter        chi.Router
	Middlewares       []MiddlewareFunc
	TaggedMiddlewares map[string]MiddlewareFunc
	ErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
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

// HandlerWithOptions creates http.Handler with additional options
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

	if options.BaseURL == "" {
		options.BaseURL = "/"
	}

	wrapper := ServerInterfaceWrapper{
		Handler: si, HandlerMiddlewares: options.Middlewares,
		TaggedMiddlewares: options.TaggedMiddlewares,
		ErrorHandlerFunc:  options.ErrorHandlerFunc,
	}

	r.Route(options.BaseURL, func(r chi.Router) {
		r.Get("/every-type-optional", wrapper.GetEveryTypeOptional)
		r.Get("/get-simple", wrapper.GetSimple)
		r.Get("/get-with-args", wrapper.GetWithArgs)
		r.Get("/get-with-references/{global_argument}/{argument}", wrapper.GetWithReferences)
		r.Get("/get-with-tagged-middleware", wrapper.GetWithTaggedMiddleware)
		r.Post("/get-with-tagged-middleware", wrapper.PostWithTaggedMiddleware)
		r.Get("/get-with-type/{content_type}", wrapper.GetWithContentType)
		r.Get("/reserved-keyword", wrapper.GetReservedKeyword)
		r.Post("/resource/{argument}", wrapper.CreateResource)
		r.Post("/resource2/{inline_argument}", wrapper.CreateResource2)
		r.Put("/resource3/{fallthrough}", wrapper.UpdateResource3)
		r.Get("/response-with-reference", wrapper.GetResponseWithReference)

	})
	return r
}
