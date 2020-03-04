package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNotFoundHandler_ServeHTTP(test *testing.T) {
	type fields struct {
		TextErrorPresenter ErrorPresenter
		JSONErrorPresenter ErrorPresenter
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
			name: "without the Accept header",
			fields: fields{
				TextErrorPresenter: new(MockErrorPresenter),
				JSONErrorPresenter: func() ErrorPresenter {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)

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
				request: httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
			},
		},
		{
			name: "with `text/html` in the Accept header",
			fields: fields{
				TextErrorPresenter: func() ErrorPresenter {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Add("Accept", "text/html")

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
				JSONErrorPresenter: new(MockErrorPresenter),
			},
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Add("Accept", "text/html")

					return request
				}(),
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			writer := httptest.NewRecorder()
			handler := NotFoundHandler{
				TextErrorPresenter: data.fields.TextErrorPresenter,
				JSONErrorPresenter: data.fields.JSONErrorPresenter,
			}
			handler.ServeHTTP(writer, data.args.request)

			response := writer.Result()
			responseBody, _ := ioutil.ReadAll(response.Body)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.TextErrorPresenter,
				data.fields.JSONErrorPresenter,
			)
			assert.Empty(test, responseBody)
		})
	}
}

func Test_isTextAccepted(test *testing.T) {
	type args struct {
		request *http.Request
	}

	for _, data := range []struct {
		name string
		args args
		want assert.BoolAssertionFunc
	}{
		{
			name: "without the Accept header",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
			},
			want: assert.False,
		},
		{
			name: "with `text/html` in the Accept header (GET)",
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Add("Accept", "text/html")

					return request
				}(),
			},
			want: assert.True,
		},
		{
			name: "with `text/html` in the Accept header (POST)",
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodPost, "http://example.com/", nil)
					request.Header.Add("Accept", "text/html")

					return request
				}(),
			},
			want: assert.False,
		},
		{
			name: "with `application/json` in the Accept header",
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Add("Accept", "application/json")

					return request
				}(),
			},
			want: assert.False,
		},
		{
			name: "with `application/json` and `text/html` in the Accept header",
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Add("Accept", "application/json")
					request.Header.Add("Accept", "text/html")

					return request
				}(),
			},
			want: assert.True,
		},
		{
			name: "with `text/*` in the Accept header",
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Add("Accept", "text/*")

					return request
				}(),
			},
			want: assert.False,
		},
		{
			name: "with `*/*` in the Accept header",
			args: args{
				request: func() *http.Request {
					request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
					request.Header.Add("Accept", "*/*")

					return request
				}(),
			},
			want: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := isTextAccepted(data.args.request)

			data.want(test, got)
		})
	}
}
