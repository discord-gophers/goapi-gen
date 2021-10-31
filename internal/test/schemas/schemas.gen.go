// Package schemas provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/discord-gophers/goapi-gen version (devel) DO NOT EDIT.
package schemas

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/discord-gophers/goapi-gen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	AccessTokenScopes = "access_token.Scopes"
)

// Defines values for EnumInObjInArrayVal.
var (
	UnknownEnumInObjInArrayVal = EnumInObjInArrayVal{}

	EnumInObjInArrayValFirst = EnumInObjInArrayVal{"first"}

	EnumInObjInArrayValSecond = EnumInObjInArrayVal{"second"}
)

// This schema name starts with a number
type N5startsWithNumber map[string]interface{}

// AnyType1 defines model for AnyType1.
type AnyType1 interface{}

// AnyType2 represents any type.
//
// This should be an interface{}
type AnyType2 interface{}

// CustomStringType defines model for CustomStringType.
type CustomStringType string

// EnumInObjInArray defines model for EnumInObjInArray.
type EnumInObjInArray []struct {
	Val *EnumInObjInArrayVal `json:"val,omitempty"`
}

// GenericObject defines model for GenericObject.
type GenericObject map[string]interface{}

// NullableProperties defines model for NullableProperties.
type NullableProperties struct {
	Optional            *string `json:"optional,omitempty"`
	OptionalAndNullable *string `json:"optionalAndNullable"`
	Required            string  `json:"required"`
	RequiredAndNullable *string `json:"requiredAndNullable"`
}

// StringInPath defines model for StringInPath.
type StringInPath string

// EnumInObjInArrayVal defines model for EnumInObjInArray.Val.
type EnumInObjInArrayVal struct {
	value string
}

func (t EnumInObjInArrayVal) ToValue() string {
	return t.value
}
func (t EnumInObjInArrayVal) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}
func (t EnumInObjInArrayVal) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return nil
}
func (t EnumInObjInArrayVal) FromValue(value string) error {
	switch value {

	case EnumInObjInArrayValFirst.value:
		t.value = value
		return nil

	case EnumInObjInArrayValSecond.value:
		t.value = value
		return nil

	}
	return fmt.Errorf("unknown enum value: %v", value)
}

// Issue185JSONBody defines parameters for Issue185.
type Issue185JSONBody NullableProperties

// Issue9JSONBody defines parameters for Issue9.
type Issue9JSONBody interface{}

// Issue9Params defines parameters for Issue9.
type Issue9Params struct {
	Foo string `json:"foo"`
}

// Issue185JSONRequestBody defines body for Issue185 for application/json ContentType.
type Issue185JSONRequestBody Issue185JSONBody

// Bind implements render.Binder.
func (Issue185JSONRequestBody) Bind(*http.Request) error {
	return nil
}

// Issue9JSONRequestBody defines body for Issue9 for application/json ContentType.
type Issue9JSONRequestBody Issue9JSONBody

// Response is a common response struct for all the API calls.
// A Response object may be instantiated via functions for specific operation responses.
type Response struct {
	body        interface{}
	statusCode  int
	contentType string
}

// Render implements the render.Renderer interface. It sets the Content-Type header
// and status code based on the response definition.
func (resp *Response) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", resp.contentType)
	render.Status(r, resp.statusCode)
	return nil
}

// Status is a builder method to override the default status code for a response.
func (resp *Response) Status(statusCode int) *Response {
	resp.statusCode = statusCode
	return resp
}

// ContentType is a builder method to override the default content type for a response.
func (resp *Response) ContentType(contentType string) *Response {
	resp.contentType = contentType
	return resp
}

// MarshalJSON implements the json.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.body)
}

// MarshalXML implements the xml.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.Encode(resp.body)
}

// EnsureEverythingIsReferencedJSON200Response is a constructor method for a EnsureEverythingIsReferenced response.
// A *Response is returned with the configured status code and content type from the spec.
func EnsureEverythingIsReferencedJSON200Response(body struct {
	AnyType1 *AnyType1 `json:"anyType1,omitempty"`

	// AnyType2 represents any type.
	//
	// This should be an interface{}
	AnyType2         *AnyType2         `json:"anyType2,omitempty"`
	CustomStringType *CustomStringType `foo:"bar" json:"customStringType,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "application/json",
	}
}

// Issue127JSON200Response is a constructor method for a Issue127 response.
// A *Response is returned with the configured status code and content type from the spec.
func Issue127JSON200Response(body GenericObject) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "application/json",
	}
}

// Issue127XML200Response is a constructor method for a Issue127 response.
// A *Response is returned with the configured status code and content type from the spec.
func Issue127XML200Response(body GenericObject) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "application/xml",
	}
}

// Issue127YAML200Response is a constructor method for a Issue127 response.
// A *Response is returned with the configured status code and content type from the spec.
func Issue127YAML200Response(body GenericObject) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "text/yaml",
	}
}

// Issue127JSONDefaultResponse is a constructor method for a Issue127 response.
// A *Response is returned with the configured status code and content type from the spec.
func Issue127JSONDefaultResponse(body GenericObject) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "application/json",
	}
}

// GetIssues375JSON200Response is a constructor method for a GetIssues375 response.
// A *Response is returned with the configured status code and content type from the spec.
func GetIssues375JSON200Response(body EnumInObjInArray) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "application/json",
	}
}

// ClientInterface is implemented by Client
type ClientInterface interface {
	// EnsureEverythingIsReferenced makes the request to the API endpoint.
	EnsureEverythingIsReferenced(ctx context.Context, opts ...func(*http.Request) error) error
	// Issue127 makes the request to the API endpoint.
	Issue127(ctx context.Context, opts ...func(*http.Request) error) error
	// Issue185 makes the request to the API endpoint.
	Issue185(ctx context.Context, respBody interface{}, params Issue185ClientParams, opts ...func(*http.Request) error) (*ReqResponse, error)
	// Issue209 makes the request to the API endpoint.
	Issue209(ctx context.Context, params Issue209ClientParams, opts ...func(*http.Request) error) error
	// Issue30 makes the request to the API endpoint.
	Issue30(ctx context.Context, params Issue30ClientParams, opts ...func(*http.Request) error) error
	// GetIssues375 makes the request to the API endpoint.
	GetIssues375(ctx context.Context, opts ...func(*http.Request) error) error
	// Issue41 makes the request to the API endpoint.
	Issue41(ctx context.Context, params Issue41ClientParams, opts ...func(*http.Request) error) error
	// Issue9 makes the request to the API endpoint.
	Issue9(ctx context.Context, respBody interface{}, params Issue9ClientParams, opts ...func(*http.Request) error) (*ReqResponse, error)
}

// Doer performs HTTP requests.
// The standard http.Client implements this interface.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://example.com for example. This can contain a path relative
	// to the server, such as https://example.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	BaseURL string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	client Doer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	reqEditors []func(req *http.Request) error
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		BaseURL: server,
		client:  &http.Client{},
	}

	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}

	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.BaseURL, "/") {
		client.BaseURL += "/"
	}

	return &client, nil
}

// WithDoer allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithDoer(doer Doer) ClientOption {
	return func(c *Client) error {
		c.client = doer
		return nil
	}
}

// WithEditors allows setting up request editors, which are used to modify
func WithEditors(fns ...func(req *http.Request) error) ClientOption {
	return func(c *Client) error {
		c.reqEditors = append(c.reqEditors, fns...)
		return nil
	}
}

type ReqResponse struct {
	*http.Response
}

// Decode is a package-level variable set to our default Decoder. We do this
// because it allows you to set Decode to another function with the
// same function signature, while also utilizing the Decoder() function
// itself. Effectively, allowing you to easily add your own logic to the package
// defaults. For example, maybe you want to impose a limit on the number of
// bytes allowed to be read from the request body.
var ReqDecoder = defaultDecoder

// defaultDecoder detects the correct decoder for use on an HTTP request and
// marshals into a given interface.
func defaultDecoder(resp *http.Response, v interface{}) error {
	var err error

	switch render.GetContentType(resp.Header.Get("Content-Type")) {
	case render.ContentTypeJSON:
		err = render.DecodeJSON(resp.Body, v)
	case render.ContentTypeXML:
		err = render.DecodeXML(resp.Body, v)
	default:
		err = errors.New("defaultDecoder: unable to automatically decode the request content type")
	}

	return err
}

// We generate a new type for each client function such that we have all required in this parameter.
// Having a parameter like this is good because we don't break the function signature if things change inside.
// This is also cleaner than having all parameters as function parameters.
// The only issue is that it easily gets quite big

type Issue185ClientParams struct {
	Body io.Reader
}

type Issue209ClientParams struct {
	// A string path parameter
	Str string `json:"str"`
}

type Issue30ClientParams struct {
	PFallthrough string `json:"fallthrough"`
}

type Issue41ClientParams struct {
	N1param string `json:"1param"`
}

type Issue9ClientParams struct {
	// Optional body
	Body io.Reader
	Foo  string `json:"foo"`
}

func buildURL(baseURL string, pathParams map[string]string, queryParams map[string]string) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	// add path parameters
	for name, value := range pathParams {
		u.Path = strings.Replace(u.Path, "{"+name+"}", value, 1)
	}

	// add query parameters
	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

// EnsureEverythingIsReferenced makes the request to the API endpoint.
func (c *Client) EnsureEverythingIsReferenced(ctx context.Context, opts ...func(*http.Request) error) error {

	// Create the request
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.BaseURL,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	// Apply any request editors
	for _, fn := range c.reqEditors {
		if err := fn(req); err != nil {
			return fmt.Errorf("failed to apply request editor: %w", err)
		}
	}

	// Do the request
	_, errDo := c.client.Do(req)
	if errDo != nil {
		return fmt.Errorf("failed to send request: %w", errDo)
	}

	return nil
}

// Issue127 makes the request to the API endpoint.
func (c *Client) Issue127(ctx context.Context, opts ...func(*http.Request) error) error {

	// Create the request
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.BaseURL,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	// Apply any request editors
	for _, fn := range c.reqEditors {
		if err := fn(req); err != nil {
			return fmt.Errorf("failed to apply request editor: %w", err)
		}
	}

	// Do the request
	_, errDo := c.client.Do(req)
	if errDo != nil {
		return fmt.Errorf("failed to send request: %w", errDo)
	}

	return nil
}

// Issue185 makes the request to the API endpoint.
func (c *Client) Issue185(ctx context.Context, respBody interface{}, params Issue185ClientParams, opts ...func(*http.Request) error) (*ReqResponse, error) {

	// Create the request
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.BaseURL,
		params.Body,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	// Apply any request editors
	for _, fn := range c.reqEditors {
		if err := fn(req); err != nil {
			return nil, fmt.Errorf("failed to apply request editor: %w", err)
		}
	}

	// Do the request
	resp, errDo := c.client.Do(req)
	if errDo != nil {
		return &ReqResponse{
			Response: resp,
		}, fmt.Errorf("failed to send request: %w", errDo)
	}

	// Bind the response body
	if err := ReqDecoder(resp, respBody); err != nil {
		return &ReqResponse{
				Response: resp,
			},
			fmt.Errorf("failed to bind response body: %w", err)
	}

	return &ReqResponse{
		Response: resp,
	}, nil
}

// Issue209 makes the request to the API endpoint.
func (c *Client) Issue209(ctx context.Context, params Issue209ClientParams, opts ...func(*http.Request) error) error {

	// Create the request
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",

		buildURL(
			c.BaseURL,
			map[string]string{
				"str": params.Str,
			},
			nil,
		),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	// Apply any request editors
	for _, fn := range c.reqEditors {
		if err := fn(req); err != nil {
			return fmt.Errorf("failed to apply request editor: %w", err)
		}
	}

	// Do the request
	_, errDo := c.client.Do(req)
	if errDo != nil {
		return fmt.Errorf("failed to send request: %w", errDo)
	}

	return nil
}

// Issue30 makes the request to the API endpoint.
func (c *Client) Issue30(ctx context.Context, params Issue30ClientParams, opts ...func(*http.Request) error) error {

	// Create the request
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",

		buildURL(
			c.BaseURL,
			map[string]string{
				"fallthrough": params.PFallthrough,
			},
			nil,
		),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	// Apply any request editors
	for _, fn := range c.reqEditors {
		if err := fn(req); err != nil {
			return fmt.Errorf("failed to apply request editor: %w", err)
		}
	}

	// Do the request
	_, errDo := c.client.Do(req)
	if errDo != nil {
		return fmt.Errorf("failed to send request: %w", errDo)
	}

	return nil
}

// GetIssues375 makes the request to the API endpoint.
func (c *Client) GetIssues375(ctx context.Context, opts ...func(*http.Request) error) error {

	// Create the request
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.BaseURL,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	// Apply any request editors
	for _, fn := range c.reqEditors {
		if err := fn(req); err != nil {
			return fmt.Errorf("failed to apply request editor: %w", err)
		}
	}

	// Do the request
	_, errDo := c.client.Do(req)
	if errDo != nil {
		return fmt.Errorf("failed to send request: %w", errDo)
	}

	return nil
}

// Issue41 makes the request to the API endpoint.
func (c *Client) Issue41(ctx context.Context, params Issue41ClientParams, opts ...func(*http.Request) error) error {

	// Create the request
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",

		buildURL(
			c.BaseURL,
			map[string]string{
				"1param": params.N1param,
			},
			nil,
		),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	// Apply any request editors
	for _, fn := range c.reqEditors {
		if err := fn(req); err != nil {
			return fmt.Errorf("failed to apply request editor: %w", err)
		}
	}

	// Do the request
	_, errDo := c.client.Do(req)
	if errDo != nil {
		return fmt.Errorf("failed to send request: %w", errDo)
	}

	return nil
}

// Issue9 makes the request to the API endpoint.
func (c *Client) Issue9(ctx context.Context, respBody interface{}, params Issue9ClientParams, opts ...func(*http.Request) error) (*ReqResponse, error) {

	queryParams := make(map[string]string)
	queryParams["foo"] = params.Foo

	// Create the request
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",

		buildURL(
			c.BaseURL,
			nil,
			queryParams,
		),
		params.Body,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	// Apply any request editors
	for _, fn := range c.reqEditors {
		if err := fn(req); err != nil {
			return nil, fmt.Errorf("failed to apply request editor: %w", err)
		}
	}

	// Do the request
	resp, errDo := c.client.Do(req)
	if errDo != nil {
		return &ReqResponse{
			Response: resp,
		}, fmt.Errorf("failed to send request: %w", errDo)
	}

	// Bind the response body
	if err := ReqDecoder(resp, respBody); err != nil {
		return &ReqResponse{
				Response: resp,
			},
			fmt.Errorf("failed to bind response body: %w", err)
	}

	return &ReqResponse{
		Response: resp,
	}, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /ensure-everything-is-referenced)
	EnsureEverythingIsReferenced(w http.ResponseWriter, r *http.Request)

	// (GET /issues/127)
	Issue127(w http.ResponseWriter, r *http.Request)

	// (GET /issues/185)
	Issue185(w http.ResponseWriter, r *http.Request)

	// (GET /issues/209/${str})
	Issue209(w http.ResponseWriter, r *http.Request, str StringInPath)

	// (GET /issues/30/{fallthrough})
	Issue30(w http.ResponseWriter, r *http.Request, pFallthrough string)

	// (GET /issues/375)
	GetIssues375(w http.ResponseWriter, r *http.Request)

	// (GET /issues/41/{1param})
	Issue41(w http.ResponseWriter, r *http.Request, n1param N5startsWithNumber)

	// (GET /issues/9)
	Issue9(w http.ResponseWriter, r *http.Request, params Issue9Params)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler          ServerInterface
	Middlewares      map[string]func(http.Handler) http.Handler
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// EnsureEverythingIsReferenced operation middleware
func (siw *ServerInterfaceWrapper) EnsureEverythingIsReferenced(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, AccessTokenScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.EnsureEverythingIsReferenced(w, r)
	})

	handler(w, r.WithContext(ctx))
}

// Issue127 operation middleware
func (siw *ServerInterfaceWrapper) Issue127(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, AccessTokenScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Issue127(w, r)
	})

	handler(w, r.WithContext(ctx))
}

// Issue185 operation middleware
func (siw *ServerInterfaceWrapper) Issue185(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, AccessTokenScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Issue185(w, r)
	})

	handler(w, r.WithContext(ctx))
}

// Issue209 operation middleware
func (siw *ServerInterfaceWrapper) Issue209(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "str" -------------
	var str StringInPath

	if err := runtime.BindStyledParameter("simple", false, "str", chi.URLParam(r, "str"), &str); err != nil {
		err = fmt.Errorf("invalid format for parameter str: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
		return
	}

	ctx = context.WithValue(ctx, AccessTokenScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Issue209(w, r, str)
	})

	handler(w, r.WithContext(ctx))
}

// Issue30 operation middleware
func (siw *ServerInterfaceWrapper) Issue30(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "fallthrough" -------------
	var pFallthrough string

	if err := runtime.BindStyledParameter("simple", false, "fallthrough", chi.URLParam(r, "fallthrough"), &pFallthrough); err != nil {
		err = fmt.Errorf("invalid format for parameter fallthrough: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
		return
	}

	ctx = context.WithValue(ctx, AccessTokenScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Issue30(w, r, pFallthrough)
	})

	handler(w, r.WithContext(ctx))
}

// GetIssues375 operation middleware
func (siw *ServerInterfaceWrapper) GetIssues375(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, AccessTokenScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetIssues375(w, r)
	})

	handler(w, r.WithContext(ctx))
}

// Issue41 operation middleware
func (siw *ServerInterfaceWrapper) Issue41(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "1param" -------------
	var n1param N5startsWithNumber

	if err := runtime.BindStyledParameter("simple", false, "1param", chi.URLParam(r, "1param"), &n1param); err != nil {
		err = fmt.Errorf("invalid format for parameter 1param: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
		return
	}

	ctx = context.WithValue(ctx, AccessTokenScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Issue41(w, r, n1param)
	})

	handler(w, r.WithContext(ctx))
}

// Issue9 operation middleware
func (siw *ServerInterfaceWrapper) Issue9(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, AccessTokenScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params Issue9Params

	// ------------- Required query parameter "foo" -------------

	if err := runtime.BindQueryParameter("form", true, true, "foo", r.URL.Query(), &params.Foo); err != nil {
		err = fmt.Errorf("invalid format for parameter foo: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Issue9(w, r, params)
	})

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

type ServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      map[string]func(http.Handler) http.Handler
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

type ServerOption func(*ServerOptions)

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface, opts ...ServerOption) http.Handler {
	options := &ServerOptions{
		BaseURL:     "/",
		BaseRouter:  chi.NewRouter(),
		Middlewares: make(map[string]func(http.Handler) http.Handler),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}

	for _, f := range opts {
		f(options)
	}

	r := options.BaseRouter
	wrapper := ServerInterfaceWrapper{
		Handler:          si,
		Middlewares:      options.Middlewares,
		ErrorHandlerFunc: options.ErrorHandlerFunc,
	}

	r.Route(options.BaseURL, func(r chi.Router) {
		r.Get("/ensure-everything-is-referenced", wrapper.EnsureEverythingIsReferenced)
		r.Get("/issues/127", wrapper.Issue127)
		r.Get("/issues/185", wrapper.Issue185)
		r.Get("/issues/209/${str}", wrapper.Issue209)
		r.Get("/issues/30/{fallthrough}", wrapper.Issue30)
		r.Get("/issues/375", wrapper.GetIssues375)
		r.Get("/issues/41/{1param}", wrapper.Issue41)
		r.Get("/issues/9", wrapper.Issue9)

	})
	return r
}

func WithRouter(r chi.Router) ServerOption {
	return func(s *ServerOptions) {
		s.BaseRouter = r
	}
}

func WithServerBaseURL(url string) ServerOption {
	return func(s *ServerOptions) {
		s.BaseURL = url
	}
}

func WithMiddleware(key string, middleware func(http.Handler) http.Handler) ServerOption {
	return func(s *ServerOptions) {
		s.Middlewares[key] = middleware
	}
}

func WithMiddlewares(middlewares map[string]func(http.Handler) http.Handler) ServerOption {
	return func(s *ServerOptions) {
		s.Middlewares = middlewares
	}
}

func WithErrorHandler(handler func(w http.ResponseWriter, r *http.Request, err error)) ServerOption {
	return func(s *ServerOptions) {
		s.ErrorHandlerFunc = handler
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RX0W/buA/+Vwj9BvxenDjtNmzNW2/YDT3gtmItsIemD4rNxFptypPoNkaQ//1AyamT",
	"xelt160vjS1R5PeR/ESvVWar2hISezVdq1o7XSGjC09X7AwtL+hScyHPOfrMmZqNJTVV5+DDOtSaC3i0",
	"VIkysixvVaJIV6imyrMsOPzWGIe5mrJrMFE+K7DScjS3dbfN0FJtNpvtYgjk9RVrx/6L4eJjU83RHUZz",
	"XRgP0QTEJ/hgAg+GC9BA0SzZOrLzr5ix2iTqnNrrtsYTNV33T6cDcLsVcFg79MIYaGpBDhzPaEYxgsI2",
	"ZQ5zBE1giNEtdIbrzYzE17vGs60irdchkLVaWFdpVlOVhcU+xI6LRK1GSzvCFTs9Yr300ciqqZprp4Sp",
	"99RUF/Rp/vWCzp3TrewwjFVMqbM1OjYYnu51Kf+QmkpNb9TCOM8qUR4zS7m6TQ4SMcBY90IHV5tEfUBC",
	"Z7JPcUOfzN7iY1OWel7i5V4s+5HZQHQM77sgksfFc8q3Z8k+evwd6+nAri+49fHFnzt079Sb/vfwebcH",
	"/G0C240z3F5JuUb0OsvQ+xHbOyR5nqN26P7c1sZfX65HsVAg7oSwczwj1TWKuIhGfQUVzHXsJUMLO9Az",
	"6Bky7dHDwjq4187YxoPxvgmvGsrB3qMDNhWO4bJE7RF0noMG3tqK6YykE+bNEhZmhXkMiw0LidHLFbr7",
	"ENo9Oh+9n4wn40lMLpKujZqql+PJ+EQlQTsCLSmSbxyO8B5dy4Wh5cj4kcMFOqQs5nWJfEQOkPLaGmLA",
	"lfHswVvgQjP0mgeZJmnWzKFmzMEQcGH8jHyNGWjKgSzLhto1hHnAJUWrxc1FrqbqfQjw/WN8F/5zH53U",
	"hK8t+Zjk08lE/mWWGCkEreu6NFk4Lf3qbUh9L4r7DaJ7oVIvHC7UVP0v7aGknV6mj4K2SbY2pz9ocyo2",
	"2YBIPWV7IGoDohH/EpXG2kpPTt8cTd3f+g5BSIWGfFPX1klmAmkrDnLrIbf0f4baIVY1Q78rrI4H0nQh",
	"fsXrM1PyFBH7Oihwd89aVeVzjhLwaaXdXW4f6NkHtfo50cgxOS50U/JvJO8XIf6+8t6+Pi4abY2wFPuA",
	"AB4KJNhePelW3qFvS9AOYXtfHC+7t6+72wE9/2Hz9peRNnCvRrQ7NS7h7RJwOjlLX6w9u81RHt4VmN15",
	"MIt+qotQc8xK3VNQtsOATydn6jCGZG+6vBlG1m9J96bPze0OhJeTdL3QZcmFs82y2Bwi+IxeLpwc7rB9",
	"sC7fHcxqh+GWErGXK08IDCNjJxwdJQO4Xk5+BNbA9LsT7E9NwXug3xwvXBkAu+R0lav9tpBFFR9MhpJO",
	"LhBk9AvrhmRGjQo9o4fCZEX33pscwS5kOQx5Q5X9ATlw4iWu3yiqB7PtQUe/OknXJyEHxyv6cpuinW8D",
	"+XQJXweP3wYDKX8Vx5F/S3D0/2RunwJ5+H2z2dw+2cVnx5u3NEgcO9eHCxEMZdY5zLhs5XfZ5JiHia/T",
	"pEjD3OatjDwz6vEe1bSzI7R8a9C1O4Vv7c8V/H/Wye5S2mXiU6fcAZkaUsWdWTxA2J/Cb24lniAkHcTG",
	"ld1YPU1TXOmqLnGc2Urk6Z8AAAD//7XAlZxLDwAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
