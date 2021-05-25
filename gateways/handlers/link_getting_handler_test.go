package handlers

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/iotest"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
)

func TestLinkGettingHandler_ServeHTTP(test *testing.T) {
	type fields struct {
		LinkGetter     LinkGetter
		LinkPresenter  LinkPresenter
		ErrorPresenter ErrorPresenter
	}
	type args struct {
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
				LinkGetter: func() LinkGetter {
					getter := new(MockLinkGetter)
					getter.
						On("GetLink", "code").
						Return(entities.Link{Code: "code", URL: "url"}, nil)

					return getter
				}(),
				LinkPresenter: func() LinkPresenter {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request = mux.SetURLVars(request, map[string]string{"code": "code"})

					presenter := new(MockLinkPresenter)
					presenter.On(
						"PresentLink",
						mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
						request,
						entities.Link{Code: "code", URL: "url"},
					)

					return presenter
				}(),
				ErrorPresenter: new(MockErrorPresenter),
			},
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request = mux.SetURLVars(request, map[string]string{"code": "code"})

					return request
				}(),
			},
		},
		{
			name: "error with path parameter decoding",
			fields: fields{
				LinkGetter:    new(MockLinkGetter),
				LinkPresenter: new(MockLinkPresenter),
				ErrorPresenter: func() ErrorPresenter {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)

					presenter := new(MockErrorPresenter)
					presenter.On(
						"PresentError",
						mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
						request,
						http.StatusBadRequest,
						mock.MatchedBy(func(error) bool { return true }),
					)

					return presenter
				}(),
			},
			args: args{
				request: httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
			},
		},
		{
			name: "error with searching",
			fields: fields{
				LinkGetter: func() LinkGetter {
					getter := new(MockLinkGetter)
					getter.On("GetLink", "code").Return(entities.Link{}, sql.ErrNoRows)

					return getter
				}(),
				LinkPresenter: new(MockLinkPresenter),
				ErrorPresenter: func() ErrorPresenter {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request = mux.SetURLVars(request, map[string]string{"code": "code"})

					presenter := new(MockErrorPresenter)
					presenter.On(
						"PresentError",
						mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
						request,
						http.StatusNotFound,
						mock.MatchedBy(func(error) bool { return true }),
					)

					return presenter
				}(),
			},
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request = mux.SetURLVars(request, map[string]string{"code": "code"})

					return request
				}(),
			},
		},
		{
			name: "error with getting",
			fields: fields{
				LinkGetter: func() LinkGetter {
					getter := new(MockLinkGetter)
					getter.On("GetLink", "code").Return(entities.Link{}, iotest.ErrTimeout)

					return getter
				}(),
				LinkPresenter: new(MockLinkPresenter),
				ErrorPresenter: func() ErrorPresenter {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request = mux.SetURLVars(request, map[string]string{"code": "code"})

					presenter := new(MockErrorPresenter)
					presenter.On(
						"PresentError",
						mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
						request,
						http.StatusInternalServerError,
						mock.MatchedBy(func(error) bool { return true }),
					)

					return presenter
				}(),
			},
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request = mux.SetURLVars(request, map[string]string{"code": "code"})

					return request
				}(),
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			writer := httptest.NewRecorder()
			handler := LinkGettingHandler{
				LinkGetter:     data.fields.LinkGetter,
				LinkPresenter:  data.fields.LinkPresenter,
				ErrorPresenter: data.fields.ErrorPresenter,
			}
			handler.ServeHTTP(writer, data.args.request)

			response := writer.Result()
			responseBody, _ := ioutil.ReadAll(response.Body)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.LinkGetter,
				data.fields.LinkPresenter,
				data.fields.ErrorPresenter,
			)
			assert.Empty(test, responseBody)
		})
	}
}
