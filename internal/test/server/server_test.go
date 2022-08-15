package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Define the required middleware. If these are not defined, the handler
// definition will panic. However, we set the middlewares to be noops.
var noopMiddlewares = map[string]func(http.Handler) http.Handler{
	"pathMiddleware":      func(h http.Handler) http.Handler { return h },
	"operationMiddleware": func(h http.Handler) http.Handler { return h },
}

func TestParameters(t *testing.T) {
	m := ServerInterfaceMock{}

	m.CreateResource2Func = func(w http.ResponseWriter, r *http.Request, inlineArgument int, params CreateResource2Params) Responser {
		assert.Equal(t, 99, *params.InlineQueryArgument)
		assert.Equal(t, 1, inlineArgument)
		return nil
	}

	h := Handler(&m, WithMiddlewares(noopMiddlewares))

	req := httptest.NewRequest("POST", "http://example.com/resource2/1?inline_query_argument=99", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	assert.Equal(t, 1, len(m.CreateResource2Calls()))
}

func TestErrorHandlerFunc(t *testing.T) {
	m := ServerInterfaceMock{}

	h := Handler(&m,
		WithMiddlewares(noopMiddlewares),
		WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
			w.Header().Set("Content-Type", "application/json")
			var requiredParamError *RequiredParamError
			assert.True(t, errors.As(err, &requiredParamError))
		}))

	s := httptest.NewServer(h)
	defer s.Close()

	res, err := http.DefaultClient.Get(s.URL + "/get-with-args")
	assert.Nil(t, err)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
}

func TestOmitMiddlewares(t *testing.T) {
	defer func() {
		// panics are the expected outcomes
		_ = recover()
	}()
	m := ServerInterfaceMock{}
	Handler(&m)

	t.Error("Expected panic without providing middlewares")
}

func TestMiddlewareCalled(t *testing.T) {
	m := ServerInterfaceMock{}
	m.GetWithTaggedMiddlewareFunc = func(w http.ResponseWriter, r *http.Request) Responser { return nil }

	called := false
	mw := map[string]func(http.Handler) http.Handler{
		"pathMiddleware": func(h http.Handler) http.Handler {
			called = true
			return h
		},
		"operationMiddleware": func(h http.Handler) http.Handler { return h },
	}

	h := Handler(&m, WithMiddlewares(mw))

	s := httptest.NewServer(h)
	defer s.Close()

	_, err := http.DefaultClient.Get(s.URL + "/with-tagged-middleware")
	assert.Nil(t, err)
	assert.True(t, called)
}

func TestMiddlewareCalledWithOrder(t *testing.T) {
	m := ServerInterfaceMock{}
	m.PostWithTaggedMiddlewareFunc = func(w http.ResponseWriter, r *http.Request) Responser { return nil }

	var order []string
	mw := map[string]func(http.Handler) http.Handler{
		"pathMiddleware": func(h http.Handler) http.Handler {
			t.Log("first")
			order = append(order, "first")
			return h
		},
		"operationMiddleware": func(h http.Handler) http.Handler {
			t.Log("second")
			order = append(order, "second")
			return h
		},
	}

	h := Handler(&m, WithMiddlewares(mw))

	s := httptest.NewServer(h)
	defer s.Close()

	_, err := http.DefaultClient.Post(s.URL+"/with-tagged-middleware", "", nil)
	assert.Nil(t, err)
	assert.Equal(t, []string{"first", "second"}, order)
}
