// Package client provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/discord-gophers/goapi-gen version (devel) DO NOT EDIT.
package client

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/render"
)

// Error defines model for Error.
type Error struct {
	// Error message
	Message string `json:"message"`
}

// NewPet defines model for NewPet.
type NewPet struct {
	// Name of the pet
	Name string `json:"name"`

	// Type of the pet
	Tag *string `json:"tag,omitempty"`
}

// Pet defines model for Pet.
type Pet struct {
	// Embedded struct due to allOf(#/components/schemas/NewPet)
	NewPet `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	// Unique id of the pet
	ID int64 `json:"id"`
}

// FindPetsParams defines parameters for FindPets.
type FindPetsParams struct {
	// tags to filter by
	Tags *[]string `json:"tags,omitempty"`

	// maximum number of results to return
	Limit *int32 `json:"limit,omitempty"`
}

// AddPetJSONBody defines parameters for AddPet.
type AddPetJSONBody NewPet

// AddPetJSONRequestBody defines body for AddPet for application/json ContentType.
type AddPetJSONRequestBody AddPetJSONBody

// Bind implements render.Binder.
func (AddPetJSONRequestBody) Bind(*http.Request) error {
	return nil
}

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

// FindPetsJSON200Response is a constructor method for a FindPets response.
// A *Response is returned with the configured status code and content type from the spec.
func FindPetsJSON200Response(body []Pet) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "application/json",
	}
}

// FindPetsJSONDefaultResponse is a constructor method for a FindPets response.
// A *Response is returned with the configured status code and content type from the spec.
func FindPetsJSONDefaultResponse(body Error) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "application/json",
	}
}

// AddPetJSON201Response is a constructor method for a AddPet response.
// A *Response is returned with the configured status code and content type from the spec.
func AddPetJSON201Response(body Pet) *Response {
	return &Response{
		body:        body,
		statusCode:  201,
		contentType: "application/json",
	}
}

// AddPetJSONDefaultResponse is a constructor method for a AddPet response.
// A *Response is returned with the configured status code and content type from the spec.
func AddPetJSONDefaultResponse(body Error) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "application/json",
	}
}

// DeletePetJSONDefaultResponse is a constructor method for a DeletePet response.
// A *Response is returned with the configured status code and content type from the spec.
func DeletePetJSONDefaultResponse(body Error) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "application/json",
	}
}

// FindPetByIDJSON200Response is a constructor method for a FindPetByID response.
// A *Response is returned with the configured status code and content type from the spec.
func FindPetByIDJSON200Response(body Pet) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "application/json",
	}
}

// FindPetByIDJSONDefaultResponse is a constructor method for a FindPetByID response.
// A *Response is returned with the configured status code and content type from the spec.
func FindPetByIDJSONDefaultResponse(body Error) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "application/json",
	}
}

// ClientInterface is implemented by Client
type ClientInterface interface {
	// FindPets makes the request to the API endpoint.
	FindPets(ctx context.Context, params FindPetsClientParams, opts ...func(*http.Request) error) error
	// AddPet makes the request to the API endpoint.
	AddPet(ctx context.Context, respBody interface{}, params AddPetClientParams, opts ...func(*http.Request) error) (*ReqResponse, error)
	// DeletePet makes the request to the API endpoint.
	DeletePet(ctx context.Context, params DeletePetClientParams, opts ...func(*http.Request) error) error
	// FindPetByID makes the request to the API endpoint.
	FindPetByID(ctx context.Context, params FindPetByIDClientParams, opts ...func(*http.Request) error) error
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

// Returns all pets
type FindPetsClientParams struct {
	// tags to filter by
	Tags *string `json:"tags,omitempty"`
	// maximum number of results to return
	Limit *string `json:"limit,omitempty"`
}

// Creates a new pet
type AddPetClientParams struct {
	// Pet to add to the store
	Body io.Reader
}

// Deletes a pet by ID
type DeletePetClientParams struct {
	// ID of pet to delete
	ID string `json:"id"`
}

// Returns a pet by ID
type FindPetByIDClientParams struct {
	// ID of pet to fetch
	ID string `json:"id"`
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

// FindPets makes the request to the API endpoint.
func (c *Client) FindPets(ctx context.Context, params FindPetsClientParams, opts ...func(*http.Request) error) error {

	queryParams := make(map[string]string)
	if params.Tags != nil {
		queryParams["tags"] = *params.Tags
	}
	if params.Limit != nil {
		queryParams["limit"] = *params.Limit
	}

	// Create the request
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",

		buildURL(
			c.BaseURL,
			nil,
			queryParams,
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

// AddPet makes the request to the API endpoint.
func (c *Client) AddPet(ctx context.Context, respBody interface{}, params AddPetClientParams, opts ...func(*http.Request) error) (*ReqResponse, error) {

	// Create the request
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
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

// DeletePet makes the request to the API endpoint.
func (c *Client) DeletePet(ctx context.Context, params DeletePetClientParams, opts ...func(*http.Request) error) error {

	// Create the request
	req, err := http.NewRequestWithContext(
		ctx,
		"DELETE",

		buildURL(
			c.BaseURL,
			map[string]string{
				"id": params.ID,
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

// FindPetByID makes the request to the API endpoint.
func (c *Client) FindPetByID(ctx context.Context, params FindPetByIDClientParams, opts ...func(*http.Request) error) error {

	// Create the request
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",

		buildURL(
			c.BaseURL,
			map[string]string{
				"id": params.ID,
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