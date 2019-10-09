package router

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewRouter(test *testing.T) {
	type args struct {
		handlers Handlers
		request  *http.Request
	}

	for _, data := range []struct {
		name string
		args args
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			writer := httptest.NewRecorder()
			router := NewRouter(data.args.handlers)
			router.ServeHTTP(writer, data.args.request)

			response := writer.Result()
			responseBody, _ := ioutil.ReadAll(response.Body)

			mock.AssertExpectationsForObjects(
				test,
				data.args.handlers.LinkGettingHandler,
				data.args.handlers.LinkCreatingHandler,
				data.args.handlers.NotFoundHandler,
			)
			assert.Empty(test, string(responseBody))
		})
	}
}
