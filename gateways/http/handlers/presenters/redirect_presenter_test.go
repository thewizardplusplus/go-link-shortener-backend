package presenters

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/iotest"

	"github.com/go-log/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
)

func TestRedirectPresenter_PresentLink(test *testing.T) {
	type fields struct {
		ErrorURL string
		Logger   log.Logger
	}
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
		link    entities.Link
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
		check   func(test *testing.T, writer http.ResponseWriter)
	}{
		{
			name: "success",
			fields: fields{
				ErrorURL: "/error",
				Logger:   new(MockLogger),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/redirect/code",
					nil,
				),
				link: entities.Link{Code: "code", URL: "https://www.google.com/"},
			},
			wantErr: assert.NoError,
			check: func(test *testing.T, writer http.ResponseWriter) {
				response := writer.(*httptest.ResponseRecorder).Result()

				assert.Equal(test, http.StatusMovedPermanently, response.StatusCode)
				assert.Equal(
					test,
					"https://www.google.com/",
					response.Header.Get("Location"),
				)
			},
		},
		{
			name: "error",
			fields: fields{
				ErrorURL: "/error",
				Logger:   new(MockLogger),
			},
			args: args{
				writer: NewTimeoutResponseRecorder(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/redirect/code",
					nil,
				),
				link: entities.Link{Code: "code", URL: "https://www.google.com/"},
			},
			wantErr: assert.Error,
			check: func(test *testing.T, writer http.ResponseWriter) {
				response := writer.(TimeoutResponseRecorder).Result()

				assert.Equal(test, http.StatusMovedPermanently, response.StatusCode)
				assert.Equal(
					test,
					"https://www.google.com/",
					response.Header.Get("Location"),
				)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			presenter := RedirectPresenter{
				ErrorURL: data.fields.ErrorURL,
				Logger:   data.fields.Logger,
			}
			gotErr :=
				presenter.PresentLink(data.args.writer, data.args.request, data.args.link)

			mock.AssertExpectationsForObjects(test, data.fields.Logger)
			data.wantErr(test, gotErr)
			data.check(test, data.args.writer)
		})
	}
}

func TestRedirectPresenter_PresentError(test *testing.T) {
	type fields struct {
		ErrorURL string
		Logger   log.Logger
	}
	type args struct {
		writer     http.ResponseWriter
		request    *http.Request
		statusCode int
		err        error
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
		check   func(test *testing.T, writer http.ResponseWriter)
	}{
		{
			name: "success",
			fields: fields{
				ErrorURL: "/error",
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
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/redirect/code",
					nil,
				),
				statusCode: http.StatusInternalServerError,
				err:        iotest.ErrTimeout,
			},
			wantErr: assert.NoError,
			check: func(test *testing.T, writer http.ResponseWriter) {
				response := writer.(*httptest.ResponseRecorder).Result()

				assert.Equal(test, http.StatusFound, response.StatusCode)
				assert.Equal(test, "/error", response.Header.Get("Location"))
			},
		},
		{
			name: "error",
			fields: fields{
				ErrorURL: "/error",
				Logger:   new(MockLogger),
			},
			args: args{
				writer: NewTimeoutResponseRecorder(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/redirect/code",
					nil,
				),
				statusCode: http.StatusInternalServerError,
				err:        iotest.ErrTimeout,
			},
			wantErr: assert.Error,
			check: func(test *testing.T, writer http.ResponseWriter) {
				response := writer.(TimeoutResponseRecorder).Result()

				assert.Equal(test, http.StatusFound, response.StatusCode)
				assert.Equal(test, "/error", response.Header.Get("Location"))
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			presenter := RedirectPresenter{
				ErrorURL: data.fields.ErrorURL,
				Logger:   data.fields.Logger,
			}
			gotErr := presenter.PresentError(
				data.args.writer,
				data.args.request,
				data.args.statusCode,
				data.args.err,
			)

			mock.AssertExpectationsForObjects(test, data.fields.Logger)
			data.wantErr(test, gotErr)
			data.check(test, data.args.writer)
		})
	}
}
