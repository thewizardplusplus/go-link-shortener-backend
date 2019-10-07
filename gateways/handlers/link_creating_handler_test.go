package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLinkCreatingHandler_ServeHTTP(test *testing.T) {
	type fields struct {
		LinkCreator    LinkCreator
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
			handler := LinkCreatingHandler{
				LinkCreator:    data.fields.LinkCreator,
				LinkPresenter:  data.fields.LinkPresenter,
				ErrorPresenter: data.fields.ErrorPresenter,
			}
			handler.ServeHTTP(writer, data.args.request)

			response := writer.Result()
			responseBody, _ := ioutil.ReadAll(response.Body)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.LinkCreator,
				data.fields.LinkPresenter,
				data.fields.ErrorPresenter,
			)
			assert.Empty(test, responseBody)
		})
	}
}
