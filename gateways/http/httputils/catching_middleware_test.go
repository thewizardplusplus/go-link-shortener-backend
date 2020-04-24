package httputils

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/iotest"

	"github.com/go-log/log"
	"github.com/stretchr/testify/mock"
)

func TestCatchingMiddleware(test *testing.T) {
	type args struct {
		logger log.Logger
	}
	type middlewareArgs struct {
		next http.Handler
	}
	type handlerArgs struct {
		writer  http.ResponseWriter
		request *http.Request
	}

	for _, data := range []struct {
		name           string
		args           args
		middlewareArgs middlewareArgs
		handlerArgs    handlerArgs
	}{
		{
			name: "success",
			args: args{
				logger: new(MockLogger),
			},
			middlewareArgs: middlewareArgs{
				next: func() http.Handler {
					handler := new(MockHandler)
					handler.
						On(
							"ServeHTTP",
							mock.MatchedBy(func(writer http.ResponseWriter) bool {
								writer.Write([]byte("test"))
								return true
							}),
							httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
						).
						Return()

					return handler
				}(),
			},
			handlerArgs: handlerArgs{
				writer: func() http.ResponseWriter {
					writer := new(MockResponseWriter)
					writer.On("Write", []byte("test")).Return(4, nil)

					return writer
				}(),
				request: httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
			},
		},
		{
			name: "error",
			args: args{
				logger: func() log.Logger {
					logger := new(MockLogger)
					logger.
						On(
							"Logf",
							mock.MatchedBy(func(string) bool { return true }),
							iotest.ErrTimeout,
						).
						Return()

					return logger
				}(),
			},
			middlewareArgs: middlewareArgs{
				next: func() http.Handler {
					handler := new(MockHandler)
					handler.
						On(
							"ServeHTTP",
							mock.MatchedBy(func(writer http.ResponseWriter) bool {
								writer.Write([]byte("test"))
								return true
							}),
							httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
						).
						Return()

					return handler
				}(),
			},
			handlerArgs: handlerArgs{
				writer: func() http.ResponseWriter {
					writer := new(MockResponseWriter)
					writer.On("Write", []byte("test")).Return(2, iotest.ErrTimeout)

					return writer
				}(),
				request: httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			middleware := CatchingMiddleware(data.args.logger)
			handler := middleware(data.middlewareArgs.next)
			handler.ServeHTTP(data.handlerArgs.writer, data.handlerArgs.request)

			mock.AssertExpectationsForObjects(
				test,
				data.args.logger,
				data.middlewareArgs.next,
				data.handlerArgs.writer,
			)
		})
	}
}
