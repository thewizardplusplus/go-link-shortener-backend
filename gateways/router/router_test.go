package router

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewRouter(test *testing.T) {
	type args struct {
		handlers Handlers
		request  *http.Request
	}

	for _, data := range []struct {
		name string
		args args
	}{
		{
			name: "getting",
			args: args{
				handlers: Handlers{
					LinkGettingHandler: func() http.Handler {
						handler := new(MockHandler)
						handler.On(
							"ServeHTTP",
							mock.MatchedBy(func(writer http.ResponseWriter) bool { return true }),
							mock.MatchedBy(func(request *http.Request) bool { return true }),
						)

						return handler
					}(),
					LinkCreatingHandler: func() http.Handler {
						handler := new(MockHandler)
						return handler
					}(),
					NotFoundHandler: func() http.Handler {
						handler := new(MockHandler)
						return handler
					}(),
				},
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/api/v1/links/code",
					nil,
				),
			},
		},
		{
			name: "creating",
			args: args{
				handlers: Handlers{
					LinkGettingHandler: func() http.Handler {
						handler := new(MockHandler)
						return handler
					}(),
					LinkCreatingHandler: func() http.Handler {
						handler := new(MockHandler)
						handler.On(
							"ServeHTTP",
							mock.MatchedBy(func(writer http.ResponseWriter) bool { return true }),
							mock.MatchedBy(func(request *http.Request) bool { return true }),
						)

						return handler
					}(),
					NotFoundHandler: func() http.Handler {
						handler := new(MockHandler)
						return handler
					}(),
				},
				request: httptest.NewRequest(
					http.MethodPost,
					"http://example.com/api/v1/links/",
					nil,
				),
			},
		},
		{
			name: "incorrect method",
			args: args{
				handlers: Handlers{
					LinkGettingHandler: func() http.Handler {
						handler := new(MockHandler)
						return handler
					}(),
					LinkCreatingHandler: func() http.Handler {
						handler := new(MockHandler)
						return handler
					}(),
					NotFoundHandler: func() http.Handler {
						handler := new(MockHandler)
						handler.On(
							"ServeHTTP",
							mock.MatchedBy(func(writer http.ResponseWriter) bool { return true }),
							mock.MatchedBy(func(request *http.Request) bool { return true }),
						)

						return handler
					}(),
				},
				request: httptest.NewRequest(
					http.MethodPost,
					"http://example.com/api/v1/links/code",
					nil,
				),
			},
		},
		{
			name: "unknown endpoint",
			args: args{
				handlers: Handlers{
					LinkGettingHandler: func() http.Handler {
						handler := new(MockHandler)
						return handler
					}(),
					LinkCreatingHandler: func() http.Handler {
						handler := new(MockHandler)
						return handler
					}(),
					NotFoundHandler: func() http.Handler {
						handler := new(MockHandler)
						handler.On(
							"ServeHTTP",
							mock.MatchedBy(func(writer http.ResponseWriter) bool { return true }),
							mock.MatchedBy(func(request *http.Request) bool { return true }),
						)

						return handler
					}(),
				},
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/api/v1/incorrect",
					nil,
				),
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			writer := httptest.NewRecorder()
			router := NewRouter(data.args.handlers)
			router.ServeHTTP(writer, data.args.request)

			response := writer.Result()
			responseBody, _ := ioutil.ReadAll(response.Body)

			mock.AssertExpectationsForObjects(
				test,
				data.args.handlers.LinkGettingHandler,
				data.args.handlers.LinkCreatingHandler,
				data.args.handlers.NotFoundHandler,
			)
			assert.Empty(test, string(responseBody))
		})
	}
}
