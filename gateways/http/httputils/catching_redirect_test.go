package httputils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCatchingRedirect(test *testing.T) {
	type args struct {
		writer     http.ResponseWriter
		request    *http.Request
		url        string
		statusCode int
	}

	for _, data := range []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success",
			args: args{
				writer: func() http.ResponseWriter {
					body := fmt.Sprintf(
						"<a href=\"http://example.com/two\">%s</a>.\n\n",
						http.StatusText(http.StatusMovedPermanently),
					)

					writer := new(MockResponseWriter)
					writer.On("Header").Return(http.Header{})
					writer.On("WriteHeader", http.StatusMovedPermanently).Return()
					writer.On("Write", []byte(body)).Return(len(body), nil)

					return writer
				}(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/one",
					nil,
				),
				url:        "http://example.com/two",
				statusCode: http.StatusMovedPermanently,
			},
			wantErr: nil,
		},
		{
			name: "error",
			args: args{
				writer: func() http.ResponseWriter {
					body := fmt.Sprintf(
						"<a href=\"http://example.com/two\">%s</a>.\n\n",
						http.StatusText(http.StatusMovedPermanently),
					)

					writer := new(MockResponseWriter)
					writer.On("Header").Return(http.Header{})
					writer.On("WriteHeader", http.StatusMovedPermanently).Return()
					writer.On("Write", []byte(body)).Return(0, iotest.ErrTimeout)

					return writer
				}(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/one",
					nil,
				),
				url:        "http://example.com/two",
				statusCode: http.StatusMovedPermanently,
			},
			wantErr: iotest.ErrTimeout,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotErr := CatchingRedirect(
				data.args.writer,
				data.args.request,
				data.args.url,
				data.args.statusCode,
			)

			mock.AssertExpectationsForObjects(test, data.args.writer)
			assert.Equal(test, data.wantErr, gotErr)
		})
	}
}
