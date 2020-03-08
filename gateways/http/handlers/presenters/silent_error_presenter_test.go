package presenters

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/mock"
)

func TestSilentErrorPresenter_PresentError(test *testing.T) {
	type fields struct {
		ErrorPresenter ErrorPresenter
		Printer        Printer
	}
	type args struct {
		writer     http.ResponseWriter
		request    *http.Request
		statusCode int
		err        error
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success",
			fields: fields{
				ErrorPresenter: func() ErrorPresenter {
					request :=
						httptest.NewRequest(http.MethodGet, "http://example.com/code", nil)

					presenter := new(MockErrorPresenter)
					presenter.
						On(
							"PresentError",
							mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
							request,
							http.StatusInternalServerError,
							iotest.ErrTimeout,
						).
						Return(nil)

					return presenter
				}(),
				Printer: new(MockPrinter),
			},
			args: args{
				writer: new(MockResponseWriter),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/code",
					nil,
				),
				statusCode: http.StatusInternalServerError,
				err:        iotest.ErrTimeout,
			},
		},
		{
			name: "error",
			fields: fields{
				ErrorPresenter: func() ErrorPresenter {
					request :=
						httptest.NewRequest(http.MethodGet, "http://example.com/code", nil)

					presenter := new(MockErrorPresenter)
					presenter.
						On(
							"PresentError",
							mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
							request,
							http.StatusInternalServerError,
							iotest.ErrTimeout,
						).
						Return(iotest.ErrTimeout)

					return presenter
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
				writer: new(MockResponseWriter),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/code",
					nil,
				),
				statusCode: http.StatusInternalServerError,
				err:        iotest.ErrTimeout,
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			presenter := SilentErrorPresenter{
				ErrorPresenter: data.fields.ErrorPresenter,
				Printer:        data.fields.Printer,
			}
			presenter.PresentError(
				data.args.writer,
				data.args.request,
				data.args.statusCode,
				data.args.err,
			)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.ErrorPresenter,
				data.fields.Printer,
				data.args.writer,
			)
		})
	}
}
