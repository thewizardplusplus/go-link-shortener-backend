package presenters

import (
	"net/http"
)

// ErrorPresenter ...
type ErrorPresenter interface {
	PresentError(writer http.ResponseWriter, statusCode int, err error) error
}
