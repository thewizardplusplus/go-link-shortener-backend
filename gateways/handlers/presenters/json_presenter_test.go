package presenters

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
)

type TimeoutResponseRecorder struct {
	*httptest.ResponseRecorder
}

func NewTimeoutResponseRecorder() TimeoutResponseRecorder {
	responseRecorder := httptest.NewRecorder()
	return TimeoutResponseRecorder{responseRecorder}
}

func (TimeoutResponseRecorder) Write([]byte) (n int, err error) {
	return 0, iotest.ErrTimeout
}

func (TimeoutResponseRecorder) WriteString(string) (n int, err error) {
	return 0, iotest.ErrTimeout
}

func TestJSONPresenter_PresentLink(test *testing.T) {
	type fields struct {
		ServerID string
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
				ServerID: "server-id",
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/code",
					nil,
				),
				link: entities.Link{Code: "code", URL: "url"},
			},
			wantErr: assert.NoError,
			check: func(test *testing.T, writer http.ResponseWriter) {
				response := writer.(*httptest.ResponseRecorder).Result()
				responseBody, _ := ioutil.ReadAll(response.Body)

				assert.Equal(test, http.StatusOK, response.StatusCode)
				assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
				assert.Equal(
					test,
					`{"ServerID":"server-id","Code":"code","URL":"url"}`,
					string(responseBody),
				)
			},
		},
		{
			name: "error",
			fields: fields{
				ServerID: "server-id",
			},
			args: args{
				writer: NewTimeoutResponseRecorder(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/code",
					nil,
				),
				link: entities.Link{Code: "code", URL: "url"},
			},
			wantErr: assert.Error,
			check: func(test *testing.T, writer http.ResponseWriter) {
				response := writer.(TimeoutResponseRecorder).Result()
				responseBody, _ := ioutil.ReadAll(response.Body)

				assert.Equal(test, http.StatusOK, response.StatusCode)
				assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
				assert.Empty(test, responseBody)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			presenter := JSONPresenter{
				ServerID: data.fields.ServerID,
			}
			gotErr :=
				presenter.PresentLink(data.args.writer, data.args.request, data.args.link)

			data.wantErr(test, gotErr)
			data.check(test, data.args.writer)
		})
	}
}

func TestJSONPresenter_PresentError(test *testing.T) {
	type args struct {
		writer     http.ResponseWriter
		request    *http.Request
		statusCode int
		err        error
	}

	for _, data := range []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
		check   func(test *testing.T, writer http.ResponseWriter)
	}{
		{
			name: "success",
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/code",
					nil,
				),
				statusCode: http.StatusInternalServerError,
				err:        iotest.ErrTimeout,
			},
			wantErr: assert.NoError,
			check: func(test *testing.T, writer http.ResponseWriter) {
				response := writer.(*httptest.ResponseRecorder).Result()
				responseBody, _ := ioutil.ReadAll(response.Body)

				assert.Equal(test, http.StatusInternalServerError, response.StatusCode)
				assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
				assert.Equal(test, `{"Error":"timeout"}`, string(responseBody))
			},
		},
		{
			name: "error",
			args: args{
				writer: NewTimeoutResponseRecorder(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/code",
					nil,
				),
				statusCode: http.StatusInternalServerError,
				err:        iotest.ErrTimeout,
			},
			wantErr: assert.Error,
			check: func(test *testing.T, writer http.ResponseWriter) {
				response := writer.(TimeoutResponseRecorder).Result()
				responseBody, _ := ioutil.ReadAll(response.Body)

				assert.Equal(test, http.StatusInternalServerError, response.StatusCode)
				assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
				assert.Empty(test, responseBody)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			var presenter JSONPresenter
			gotErr := presenter.PresentError(
				data.args.writer,
				data.args.request,
				data.args.statusCode,
				data.args.err,
			)

			data.wantErr(test, gotErr)
			data.check(test, data.args.writer)
		})
	}
}
