package presenters

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/iotest"

	"github.com/go-log/log"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
)

func TestSilentLinkPresenter_PresentLink(test *testing.T) {
	type fields struct {
		LinkPresenter LinkPresenter
		Logger        log.Logger
	}
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
		link    entities.Link
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success",
			fields: fields{
				LinkPresenter: func() LinkPresenter {
					request :=
						httptest.NewRequest(http.MethodGet, "http://example.com/code", nil)

					presenter := new(MockLinkPresenter)
					presenter.
						On(
							"PresentLink",
							mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
							request,
							entities.Link{Code: "code", URL: "url"},
						).
						Return(nil)

					return presenter
				}(),
				Logger: new(MockLogger),
			},
			args: args{
				writer: new(MockResponseWriter),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/code",
					nil,
				),
				link: entities.Link{Code: "code", URL: "url"},
			},
		},
		{
			name: "error",
			fields: fields{
				LinkPresenter: func() LinkPresenter {
					request :=
						httptest.NewRequest(http.MethodGet, "http://example.com/code", nil)

					presenter := new(MockLinkPresenter)
					presenter.
						On(
							"PresentLink",
							mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
							request,
							entities.Link{Code: "code", URL: "url"},
						).
						Return(iotest.ErrTimeout)

					return presenter
				}(),
				Logger: func() log.Logger {
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
			args: args{
				writer: new(MockResponseWriter),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/code",
					nil,
				),
				link: entities.Link{Code: "code", URL: "url"},
			},
		},
	} {
		test.Run(data.name, func(t *testing.T) {
			presenter := SilentLinkPresenter{
				LinkPresenter: data.fields.LinkPresenter,
				Logger:        data.fields.Logger,
			}
			presenter.PresentLink(data.args.writer, data.args.request, data.args.link)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.LinkPresenter,
				data.fields.Logger,
				data.args.writer,
			)
		})
	}
}
