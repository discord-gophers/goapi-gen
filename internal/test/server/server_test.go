package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParameters(t *testing.T) {
	m := ServerInterfaceMock{}

	m.CreateResource2Func = func(w http.ResponseWriter, r *http.Request, inlineArgument int, params CreateResource2Params) {
		assert.Equal(t, 99, *params.InlineQueryArgument)
		assert.Equal(t, 1, inlineArgument)
	}

	h := Handler(&m)

	req := httptest.NewRequest("POST", "http://openapitest.deepmap.ai/resource2/1?inline_query_argument=99", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	assert.Equal(t, 1, len(m.CreateResource2Calls()))
}

func TestErrorHandlerFunc(t *testing.T) {
	m := ServerInterfaceMock{}

	h := HandlerWithOptions(&m, ChiServerOptions{
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			w.Header().Set("Content-Type", "application/json")
			// FIXME This should likely return a RequiredParamError, however due to binding it never does...
			//var requiredParamError *RequiredParamError
			//assert.True(t, errors.As(err, &requiredParamError))
		},
	})

	s := httptest.NewServer(h)
	defer s.Close()

	req, err := http.DefaultClient.Get(s.URL + "/get-with-args")
	assert.Nil(t, err)
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
}

func TestErrorHandlerFuncBackwardsCompatible(t *testing.T) {
	m := ServerInterfaceMock{}

	h := HandlerWithOptions(&m, ChiServerOptions{})

	s := httptest.NewServer(h)
	defer s.Close()

	req, err := http.DefaultClient.Get(s.URL + "/get-with-args")
	b, _ := ioutil.ReadAll(req.Body)
	assert.Nil(t, err)
	assert.Equal(t, "text/plain; charset=utf-8", req.Header.Get("Content-Type"))
	assert.Equal(t, "invalid format for parameter required_argument: query parameter 'required_argument' is required\n", string(b))
}
