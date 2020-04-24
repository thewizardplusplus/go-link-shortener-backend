package httputils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSPAFallbackMiddleware(test *testing.T) {
	type middlewareArgs struct {
		next http.Handler
	}
	type handlerArgs struct {
		writer  http.ResponseWriter
		request *http.Request
	}

	for _, data := range []struct {
		name           string
		middlewareArgs middlewareArgs
		handlerArgs    handlerArgs
	}{
		{
			name: "API request",
			middlewareArgs: middlewareArgs{
				next: func() http.Handler {
					handler := new(MockHandler)
					handler.
						On(
							"ServeHTTP",
							mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
							httptest.NewRequest(
								http.MethodGet,
								"http://example.com/api/v1/endpoint",
								nil,
							),
						).
						Return()

					return handler
				}(),
			},
			handlerArgs: handlerArgs{
				writer: new(MockResponseWriter),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/api/v1/endpoint",
					nil,
				),
			},
		},
		{
			name: "static asset request/on the path",
			middlewareArgs: middlewareArgs{
				next: func() http.Handler {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Set("Accept", "text/html")
					request.RequestURI = "http://example.com/path"

					handler := new(MockHandler)
					handler.
						On(
							"ServeHTTP",
							mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
							request,
						).
						Return()

					return handler
				}(),
			},
			handlerArgs: handlerArgs{
				writer: new(MockResponseWriter),
				request: func() *http.Request {
					request :=
						httptest.NewRequest(http.MethodGet, "http://example.com/path", nil)
					request.Header.Set("Accept", "text/html")

					return request
				}(),
			},
		},
		{
			name: "static asset request/to the root",
			middlewareArgs: middlewareArgs{
				next: func() http.Handler {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Set("Accept", "text/html")

					handler := new(MockHandler)
					handler.
						On(
							"ServeHTTP",
							mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
							request,
						).
						Return()

					return handler
				}(),
			},
			handlerArgs: handlerArgs{
				writer: new(MockResponseWriter),
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Set("Accept", "text/html")

					return request
				}(),
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			middleware := SPAFallbackMiddleware()
			handler := middleware(data.middlewareArgs.next)
			handler.ServeHTTP(data.handlerArgs.writer, data.handlerArgs.request)

			mock.AssertExpectationsForObjects(
				test,
				data.middlewareArgs.next,
				data.handlerArgs.writer,
			)
		})
	}
}

func Test_isStaticAssetRequest(test *testing.T) {
	type args struct {
		request *http.Request
	}

	for _, data := range []struct {
		name string
		args args
		want assert.BoolAssertionFunc
	}{
		{
			name: "true/ideal case",
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Set("Accept", "text/html")

					return request
				}(),
			},
			want: assert.True,
		},
		{
			name: "true/with additional values",
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Set(
						"Accept",
						"text/plain,text/html,application/xhtml+xml,application/xml",
					)

					return request
				}(),
			},
			want: assert.True,
		},
		{
			name: "true/with additional values and spaces",
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Set(
						"Accept",
						"text/plain, text/html, application/xhtml+xml, application/xml",
					)

					return request
				}(),
			},
			want: assert.True,
		},
		{
			name: "true/with additional values and q-values",
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Set(
						"Accept",
						"text/plain;q=0.9,"+
							"text/html;q=0.8,"+
							"application/xhtml+xml;q=0.7,"+
							"application/xml;q=0.6",
					)

					return request
				}(),
			},
			want: assert.True,
		},
		{
			name: "true/with additional values, q-values and spaces",
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Set(
						"Accept",
						"text/plain;q=0.9, "+
							"text/html;q=0.8, "+
							"application/xhtml+xml;q=0.7, "+
							"application/xml;q=0.6",
					)

					return request
				}(),
			},
			want: assert.True,
		},
		{
			name: "true/with additional values, q-values and more spaces",
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Set(
						"Accept",
						"text/plain; q=0.9, "+
							"text/html; q=0.8, "+
							"application/xhtml+xml; q=0.7, "+
							"application/xml; q=0.6",
					)

					return request
				}(),
			},
			want: assert.True,
		},
		{
			name: "false/with an incorrect method",
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodPost, "http://example.com/", nil)
					request.Header.Set("Accept", "text/html")

					return request
				}(),
			},
			want: assert.False,
		},
		{
			name: "false/without the required header",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
			},
			want: assert.False,
		},
		{
			name: "false/without the required value",
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Set("Accept", "text/css")

					return request
				}(),
			},
			want: assert.False,
		},
		{
			name: "false/with acceptance of text/*",
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Set("Accept", "text/*")

					return request
				}(),
			},
			want: assert.False,
		},
		{
			name: "false/with acceptance of */*",
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Set("Accept", "*/*")

					return request
				}(),
			},
			want: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := isStaticAssetRequest(data.args.request)

			data.want(test, got)
		})
	}
}
