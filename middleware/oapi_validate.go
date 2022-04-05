// Package middleware implements middleware function for go-chi or net/http,
// which validates incoming HTTP requests to make sure that they conform to the given OAPI 3.0 specification.
// When OAPI validation failes on the request, we return an HTTP/400.[refactor/middleware 7ad632e] diff: revert package doc comment deletion
package middleware

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

// Options to customize request validation, openapi3filter specified options will be passed through.
type Options struct {
	Options *openapi3filter.Options
	ErrRespContentType

	respContentTypeHeader string

	router routers.Router
}

// ErrRespContentType represents the support content-types for the response when a validation error occurs
type ErrRespContentType string

// Consts to expose supported Error Response Content-Types
const (
	ErrRespContentTypePlain ErrRespContentType = "text/plain"
	ErrRespContentTypeJSON  ErrRespContentType = "application/json"
	ErrRespContentTypeXML   ErrRespContentType = "application/xml"
)

// WithErrContentType sets the content type of the error response
func WithErrContentType(contentType ErrRespContentType) func(*Options) {
	return func(options *Options) {
		options.ErrRespContentType = contentType
		options.respContentTypeHeader = string(contentType) + "; charset=utf-8"
	}
}

// WithOptions sets the openapi3filter options for the middleware
func WithOptions(opt *openapi3filter.Options) func(*Options) {
	return func(options *Options) {
		options.Options = opt
	}
}

func OAPIValidator(swagger *openapi3.T, opts ...func(*Options)) func(next http.Handler) http.Handler {
	r, err := gorillamux.NewRouter(swagger)
	if err != nil {
		// user error
		panic("could not create router: " + err.Error())
	}

	options := Options{
		ErrRespContentType:    ErrRespContentTypePlain,
		respContentTypeHeader: string(ErrRespContentTypePlain) + "; charset=utf-8",
		router:                r,
	}

	for _, opt := range opts {
		opt(&options)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// validate request
			if statusCode, err := validateOAPIRequest(r, options); err != nil {
				w.Header().Set("Content-Type", options.respContentTypeHeader)
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.WriteHeader(statusCode)

				body := []byte(err.Error())
				switch options.ErrRespContentType {
				case ErrRespContentTypeJSON:
					body, _ = json.Marshal(err.Error())
				case ErrRespContentTypeXML:
					body, _ = xml.Marshal(err.Error())
				}

				fmt.Fprintln(w, string(body))
				return
			}

			// serve
			next.ServeHTTP(w, r)
		})
	}
}

// This function is called from the middleware above and actually does the work
// of validating a request.
func validateOAPIRequest(r *http.Request, options Options) (int, error) {
	// pain
	route, pathParams, err := options.router.FindRoute(r)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// Validate request
	reqValidation := &openapi3filter.RequestValidationInput{
		Request:    r,
		PathParams: pathParams,
		Route:      route,
		Options:    options.Options,
	}

	// Validate security before any other validation, unless options.Options.MultiError is true
	if options.Options == nil || !options.Options.MultiError {
		security := reqValidation.Route.Operation.Security
		if security == nil {
			security = &reqValidation.Route.Spec.Security
		}

		if err := openapi3filter.ValidateSecurityRequirements(r.Context(), reqValidation, *security); err != nil {
			return http.StatusUnauthorized, err
		}
	}

	// Validate the rest of the request
	if err := openapi3filter.ValidateRequest(r.Context(), reqValidation); err != nil {
		reqError := &openapi3filter.RequestError{}
		secError := &openapi3filter.SecurityRequirementsError{}
		otherErrr := &openapi3.MultiError{}

		switch {
		case errors.As(err, &reqError):
			// We've got a bad request
			// Split up the verbose error by lines and return the first one
			// openapi errors seem to be multi-line with a decent message on the first
			errorLines := strings.Split(err.Error(), "\n")
			return http.StatusBadRequest, fmt.Errorf(errorLines[0])
		case errors.As(err, &secError):
			return http.StatusUnauthorized, err
		case errors.As(err, &otherErrr):
			// This case occurs when options.Options.MultiError is true.
			return http.StatusBadRequest, err
		default:
			// Shouldn't happen too much
			return http.StatusInternalServerError, fmt.Errorf("error validating route: %s", err.Error())
		}
	}

	return http.StatusOK, nil
}
