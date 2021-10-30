// Package securityprovider contains some default securityprovider
// implementations, which can be used as a RequestEditorFn of a
// client.
package securityprovider

import (
	"fmt"
	"net/http"
)

const (
	// ErrAPIKeyInvalidIn indicates a usage of an invalid In.
	// Should be cookie, header or query
	ErrAPIKeyInvalidIn = Error("invalid 'in' specified for apiKey")
)

// Error defines error values of a security provider.
type Error string

// Error implements the error interface.
func (e Error) Error() string {
	return string(e)
}

// NewBasicAuth provides a BasicAuth, which can solve
// the BasicAuth challenge for api-calls.
func NewBasicAuth(username, password string) (*BasicAuth, error) {
	return &BasicAuth{
		username: username,
		password: password,
	}, nil
}

// BasicAuth sends a base64-encoded combination of
// username, password along with a request.
type BasicAuth struct {
	username string
	password string
}

// Intercept will attach an Authorization header to the request and ensures that
// the username, password are base64 encoded and attached to the header.
func (s *BasicAuth) Intercept(req *http.Request) error {
	req.SetBasicAuth(s.username, s.password)
	return nil
}

// NewBearerToken provides a BearerToken, which can solve
// the Bearer Auth challende for api-calls.
func NewBearerToken(token string) (*BearerToken, error) {
	return &BearerToken{
		token: token,
	}, nil
}

// BearerToken sends a token as part of an
// Authorization: Bearer header along with a request.
type BearerToken struct {
	token string
}

// Intercept will attach an Authorization header to the request
// and ensures that the bearer token is attached to the header.
func (s *BearerToken) Intercept(req *http.Request) error {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))
	return nil
}

// NewAPIKey will attach a generic apiKey for a given name
// either to a cookie, header or as a query parameter.
func NewAPIKey(in, name, apiKey string) (*APIKey, error) {
	interceptors := map[string]func(req *http.Request) error{
		"cookie": func(req *http.Request) error {
			req.AddCookie(&http.Cookie{Name: name, Value: apiKey})
			return nil
		},
		"header": func(req *http.Request) error {
			req.Header.Add(name, apiKey)
			return nil
		},
		"query": func(req *http.Request) error {
			query := req.URL.Query()
			query.Add(name, apiKey)
			req.URL.RawQuery = query.Encode()
			return nil
		},
	}

	interceptor, ok := interceptors[in]
	if !ok {
		return nil, ErrAPIKeyInvalidIn
	}

	return &APIKey{
		interceptor: interceptor,
	}, nil
}

// APIKey will attach an apiKey either to a
// cookie, header or query.
type APIKey struct {
	interceptor func(req *http.Request) error
}

// Intercept will attach a cookie, header or query param for the configured
// name and apiKey.
func (s *APIKey) Intercept(req *http.Request) error {
	return s.interceptor(req)
}
