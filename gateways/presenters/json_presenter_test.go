package presenters

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

type TimeoutResponseRecorder struct {
	*httptest.ResponseRecorder
}

func NewTimeoutResponseRecorder() TimeoutResponseRecorder {
	responseRecorder := httptest.NewRecorder()
	return TimeoutResponseRecorder{responseRecorder}
}

func (TimeoutResponseRecorder) WriteString(string) (n int, err error) {
	return 0, iotest.ErrTimeout
}

func TestJSONPresenter_PresentLink(test *testing.T) {
	type args struct {
		writer http.ResponseWriter
		link   entities.Link
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
				link:   entities.Link{Code: "code", URL: "url"},
			},
			wantErr: assert.NoError,
			check: func(test *testing.T, writer http.ResponseWriter) {
				response := writer.(*httptest.ResponseRecorder).Result()
				responseBody, _ := ioutil.ReadAll(response.Body)

				assert.Equal(test, http.StatusOK, response.StatusCode)
				assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
				assert.Equal(test, `{"Code":"code","URL":"url"}`, string(responseBody))
			},
		},
		{
			name: "error",
			args: args{
				writer: NewTimeoutResponseRecorder(),
				link:   entities.Link{Code: "code", URL: "url"},
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
			var presenter JSONPresenter
			gotErr := presenter.PresentLink(data.args.writer, data.args.link)

			data.wantErr(test, gotErr)
			data.check(test, data.args.writer)
		})
	}
}

func TestJSONPresenter_PresentError(test *testing.T) {
	type args struct {
		writer     http.ResponseWriter
		statusCode int
		err        error
	}

	for _, data := range []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
		check   func(test *testing.T, writer http.ResponseWriter)
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			var presenter JSONPresenter
			gotErr := presenter.PresentError(
				data.args.writer,
				data.args.statusCode,
				data.args.err,
			)

			data.wantErr(test, gotErr)
			data.check(test, data.args.writer)
		})
	}
}

func Test_presentData(test *testing.T) {
	type args struct {
		writer     http.ResponseWriter
		statusCode int
		data       interface{}
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
				writer:     httptest.NewRecorder(),
				statusCode: http.StatusOK,
				data:       entities.Link{Code: "code", URL: "url"},
			},
			wantErr: assert.NoError,
			check: func(test *testing.T, writer http.ResponseWriter) {
				response := writer.(*httptest.ResponseRecorder).Result()
				responseBody, _ := ioutil.ReadAll(response.Body)

				assert.Equal(test, http.StatusOK, response.StatusCode)
				assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
				assert.Equal(test, `{"Code":"code","URL":"url"}`, string(responseBody))
			},
		},
		{
			name: "error with marshalling",
			args: args{
				writer:     httptest.NewRecorder(),
				statusCode: http.StatusOK,
				data:       func() {},
			},
			wantErr: assert.Error,
			check: func(test *testing.T, writer http.ResponseWriter) {
				response := writer.(*httptest.ResponseRecorder).Result()
				responseBody, _ := ioutil.ReadAll(response.Body)

				assert.Equal(test, http.StatusOK, response.StatusCode)
				assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
				assert.Empty(test, responseBody)
			},
		},
		{
			name: "error with writing",
			args: args{
				writer:     NewTimeoutResponseRecorder(),
				statusCode: http.StatusOK,
				data:       entities.Link{Code: "code", URL: "url"},
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
			gotErr := presentData(data.args.writer, data.args.statusCode, data.args.data)

			data.wantErr(test, gotErr)
			data.check(test, data.args.writer)
		})
	}
}
