package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		// TODO: add test cases
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
