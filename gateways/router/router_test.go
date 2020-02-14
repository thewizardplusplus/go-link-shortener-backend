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
		redirectEndpointPrefix string
		handlers               Handlers
		request                *http.Request
	}

	for _, data := range []struct {
		name string
		args args
	}{
		{
			name: "redirect",
			args: args{
				redirectEndpointPrefix: "/redirect",
				handlers: Handlers{
					LinkRedirectHandler: func() http.Handler {
						handler := new(MockHandler)
						handler.On(
							"ServeHTTP",
							mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
							mock.MatchedBy(func(*http.Request) bool { return true }),
						)

						return handler
					}(),
					LinkGettingHandler:  new(MockHandler),
					LinkCreatingHandler: new(MockHandler),
					NotFoundHandler:     new(MockHandler),
				},
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/redirect/code",
					nil,
				),
			},
		},
		{
			name: "getting",
			args: args{
				redirectEndpointPrefix: "/redirect",
				handlers: Handlers{
					LinkRedirectHandler: new(MockHandler),
					LinkGettingHandler: func() http.Handler {
						handler := new(MockHandler)
						handler.On(
							"ServeHTTP",
							mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
							mock.MatchedBy(func(*http.Request) bool { return true }),
						)

						return handler
					}(),
					LinkCreatingHandler: new(MockHandler),
					NotFoundHandler:     new(MockHandler),
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
				redirectEndpointPrefix: "/redirect",
				handlers: Handlers{
					LinkRedirectHandler: new(MockHandler),
					LinkGettingHandler:  new(MockHandler),
					LinkCreatingHandler: func() http.Handler {
						handler := new(MockHandler)
						handler.On(
							"ServeHTTP",
							mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
							mock.MatchedBy(func(*http.Request) bool { return true }),
						)

						return handler
					}(),
					NotFoundHandler: new(MockHandler),
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
				redirectEndpointPrefix: "/redirect",
				handlers: Handlers{
					LinkRedirectHandler: new(MockHandler),
					LinkGettingHandler:  new(MockHandler),
					LinkCreatingHandler: new(MockHandler),
					NotFoundHandler: func() http.Handler {
						handler := new(MockHandler)
						handler.On(
							"ServeHTTP",
							mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
							mock.MatchedBy(func(*http.Request) bool { return true }),
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
				redirectEndpointPrefix: "/redirect",
				handlers: Handlers{
					LinkRedirectHandler: new(MockHandler),
					LinkGettingHandler:  new(MockHandler),
					LinkCreatingHandler: new(MockHandler),
					NotFoundHandler: func() http.Handler {
						handler := new(MockHandler)
						handler.On(
							"ServeHTTP",
							mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
							mock.MatchedBy(func(*http.Request) bool { return true }),
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
			router := NewRouter(data.args.redirectEndpointPrefix, data.args.handlers)
			router.ServeHTTP(writer, data.args.request)

			response := writer.Result()
			responseBody, _ := ioutil.ReadAll(response.Body)

			mock.AssertExpectationsForObjects(
				test,
				data.args.handlers.LinkRedirectHandler,
				data.args.handlers.LinkGettingHandler,
				data.args.handlers.LinkCreatingHandler,
				data.args.handlers.NotFoundHandler,
			)
			assert.Empty(test, string(responseBody))
		})
	}
}
