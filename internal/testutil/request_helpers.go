// Copyright 2019 DeepMap, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutil

// This is a set of fluent request builders for tests, which help us to
// simplify constructing and unmarshaling test objects. For example, to post
// a body and return a response, you would do something like:
//
//   var body RequestBody
//   var response ResponseBody
//   t is *testing.T, from a unit test
//   h is http.Handler
//   response := NewRequest().Post("/path").WithJsonBody(body).GoWithHTTPHandler(t, h)
//   err := response.UnmarshalBodyToObject(&response)
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// NewRequest makes a new request.
func NewRequest() *RequestBuilder {
	return &RequestBuilder{
		Headers: make(map[string]string),
	}
}

// RequestBuilder caches request settings as we build up the request.
type RequestBuilder struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    []byte
	Error   error
	Cookies []*http.Cookie
}

// WithMethod sets the method and path for an operation.
func (r *RequestBuilder) WithMethod(method string, path string) *RequestBuilder {
	r.Method = method
	r.Path = path
	return r
}

// Get does the request with the GET Method.
func (r *RequestBuilder) Get(path string) *RequestBuilder {
	return r.WithMethod("GET", path)
}

// Post does the request with the POST Method.
func (r *RequestBuilder) Post(path string) *RequestBuilder {
	return r.WithMethod("POST", path)
}

// Put does the request with the PUT method.
func (r *RequestBuilder) Put(path string) *RequestBuilder {
	return r.WithMethod("PUT", path)
}

// Patch does the request method with the PATCH method.
func (r *RequestBuilder) Patch(path string) *RequestBuilder {
	return r.WithMethod("PATCH", path)
}

// Delete does the request with the DELETE method.
func (r *RequestBuilder) Delete(path string) *RequestBuilder {
	return r.WithMethod("DELETE", path)
}

// WithHeader sets the header on r.
func (r *RequestBuilder) WithHeader(header, value string) *RequestBuilder {
	r.Headers[header] = value
	return r
}

// WithJWSAuth sets the authorization token on r.
func (r *RequestBuilder) WithJWSAuth(jws string) *RequestBuilder {
	r.Headers["Authorization"] = "Bearer " + jws
	return r
}

// WithHost sets the Host header on r.
func (r *RequestBuilder) WithHost(value string) *RequestBuilder {
	return r.WithHeader("Host", value)
}

// WithContentType sets the Content Type on r.
func (r *RequestBuilder) WithContentType(value string) *RequestBuilder {
	return r.WithHeader("Content-Type", value)
}

// WithJSONContentType sets the content type to json on r.
func (r *RequestBuilder) WithJSONContentType() *RequestBuilder {
	return r.WithContentType("application/json")
}

// WithAccept sets teh accept header on r.
func (r *RequestBuilder) WithAccept(value string) *RequestBuilder {
	return r.WithHeader("Accept", value)
}

// WithAcceptJSON sets the accept header to accept json on r.
func (r *RequestBuilder) WithAcceptJSON() *RequestBuilder {
	return r.WithAccept("application/json")
}

// Request body operations

// WithBody sets the body on r.
func (r *RequestBuilder) WithBody(body []byte) *RequestBuilder {
	r.Body = body
	return r
}

// WithJSONBody sets the body content and marshals obj as JSON.
func (r *RequestBuilder) WithJSONBody(obj interface{}) *RequestBuilder {
	var err error
	r.Body, err = json.Marshal(obj)
	if err != nil {
		r.Error = fmt.Errorf("failed to marshal json object: %s", err)
	}
	return r.WithJSONContentType()
}

// Cookie operations

// WithCookie sets the cookie on r.
func (r *RequestBuilder) WithCookie(c *http.Cookie) *RequestBuilder {
	r.Cookies = append(r.Cookies, c)
	return r
}

// WithCookieNameValue builds a cookie and adds it to r.
func (r *RequestBuilder) WithCookieNameValue(name, value string) *RequestBuilder {
	return r.WithCookie(&http.Cookie{Name: name, Value: value})
}

// GoWithHTTPHandler performs the request, it takes a pointer to a testing context
// to print messages, and a http handler for request handling.
func (r *RequestBuilder) GoWithHTTPHandler(t *testing.T, handler http.Handler) *CompletedRequest {
	if r.Error != nil {
		// Fail the test if we had an error
		t.Errorf("error constructing request: %s", r.Error)
		return nil
	}
	var bodyReader io.Reader
	if r.Body != nil {
		bodyReader = bytes.NewReader(r.Body)
	}

	req := httptest.NewRequest(r.Method, r.Path, bodyReader)
	for h, v := range r.Headers {
		req.Header.Add(h, v)
	}
	if host, ok := r.Headers["Host"]; ok {
		req.Host = host
	}
	for _, c := range r.Cookies {
		req.AddCookie(c)
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	return &CompletedRequest{
		Recorder: rec,
	}
}

// CompletedRequest represents the returned request with some helper functions.
type CompletedRequest struct {
	Recorder *httptest.ResponseRecorder

	// When set to true, decoders will be more strict. In the default JSON
	// recorder, unknown fields will cause errors.
	Strict bool
}

// DisallowUnknownFields makes the unmarshaler strict.
func (c *CompletedRequest) DisallowUnknownFields() {
	c.Strict = true
}

// UnmarshalBodyToObject unmarshales the response body.
func (c *CompletedRequest) UnmarshalBodyToObject(obj interface{}) error {
	ctype := c.Recorder.Header().Get("Content-Type")

	// Content type can have an annotation after ;
	contentParts := strings.Split(ctype, ";")
	content := strings.TrimSpace(contentParts[0])
	handler := getHandler(content)
	if handler == nil {
		return fmt.Errorf("unhandled content: %s", content)
	}

	return handler(ctype, c.Recorder.Body, obj, c.Strict)
}

// UnmarshalJSONToObject unmarshals the resposne body.
func (c *CompletedRequest) UnmarshalJSONToObject(obj interface{}) error {
	return json.Unmarshal(c.Recorder.Body.Bytes(), obj)
}

// Code returns the response code.
func (c *CompletedRequest) Code() int {
	return c.Recorder.Code
}
