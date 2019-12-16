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

func TestJSONPresenter_PresentLink(test *testing.T) {
	writer := httptest.NewRecorder()
	presenter := JSONPresenter{}
	presenter.PresentLink(writer, entities.Link{Code: "code", URL: "url"})

	response := writer.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)

	assert.Equal(test, http.StatusOK, response.StatusCode)
	assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(test, `{"Code":"code","URL":"url"}`, string(responseBody))
}

func TestJSONPresenter_PresentError(test *testing.T) {
	writer := httptest.NewRecorder()
	presenter := JSONPresenter{}
	presenter.
		PresentError(writer, http.StatusInternalServerError, iotest.ErrTimeout)

	response := writer.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)

	assert.Equal(test, http.StatusInternalServerError, response.StatusCode)
	assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(test, `{"Error":"timeout"}`, string(responseBody))
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
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			gotErr := presentData(data.args.writer, data.args.statusCode, data.args.data)

			data.wantErr(test, gotErr)
			data.check(test, data.args.writer)
		})
	}
}
