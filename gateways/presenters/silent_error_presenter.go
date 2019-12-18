package presenters

import (
	"net/http"
)

//go:generate mockery -name=ErrorPresenter -inpkg -case=underscore -testonly

// ErrorPresenter ...
type ErrorPresenter interface {
	PresentError(writer http.ResponseWriter, statusCode int, err error) error
}
