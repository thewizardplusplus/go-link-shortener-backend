package httputils

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/mock"
)

func TestCatchingHandler_ServeHTTP(test *testing.T) {
	type fields struct {
		Handler http.Handler
		Printer Printer
	}
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success",
			fields: fields{
				Handler: func() http.Handler {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)

					handler := new(MockHandler)
					handler.
						On(
							"ServeHTTP",
							mock.MatchedBy(func(writer http.ResponseWriter) bool {
								writer.Write([]byte("test"))
								return true
							}),
							request,
						).
						Return()

					return handler
				}(),
				Printer: new(MockPrinter),
			},
			args: args{
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
			fields: fields{
				Handler: func() http.Handler {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)

					handler := new(MockHandler)
					handler.
						On(
							"ServeHTTP",
							mock.MatchedBy(func(writer http.ResponseWriter) bool {
								writer.Write([]byte("test"))
								return true
							}),
							request,
						).
						Return()

					return handler
				}(),
				Printer: func() Printer {
					printer := new(MockPrinter)
					printer.
						On(
							"Printf",
							mock.MatchedBy(func(string) bool { return true }),
							iotest.ErrTimeout,
						).
						Return()

					return printer
				}(),
			},
			args: args{
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
			handler := CatchingHandler{
				Handler: data.fields.Handler,
				Printer: data.fields.Printer,
			}
			handler.ServeHTTP(data.args.writer, data.args.request)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.Handler,
				data.fields.Printer,
				data.args.writer,
			)
		})
	}
}
