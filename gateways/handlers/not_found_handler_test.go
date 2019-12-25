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
	request := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)

	presenter := new(MockErrorPresenter)
	presenter.On(
		"PresentError",
		mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
		request,
		http.StatusNotFound,
		mock.MatchedBy(func(error) bool { return true }),
	)

	writer := httptest.NewRecorder()
	handler := NotFoundHandler{
		ErrorPresenter: presenter,
	}
	handler.ServeHTTP(writer, request)

	response := writer.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)

	mock.AssertExpectationsForObjects(test, presenter)
	assert.Empty(test, responseBody)
}
