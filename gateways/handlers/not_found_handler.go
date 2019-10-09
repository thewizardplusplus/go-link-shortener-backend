package handlers

import (
	"net/http"

	"github.com/pkg/errors"
)

// NotFoundHandler ...
type NotFoundHandler struct {
	ErrorPresenter ErrorPresenter
}

// ServeHTTP ...
func (handler NotFoundHandler) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
) {
	const statusCode = http.StatusNotFound
	err := errors.New("unable to find the endpoint")
	handler.ErrorPresenter.PresentError(writer, statusCode, err)
}
