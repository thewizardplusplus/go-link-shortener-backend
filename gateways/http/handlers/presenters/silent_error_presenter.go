package presenters

import (
	"net/http"

	"github.com/go-log/log"
)

//go:generate mockery -name=ErrorPresenter -inpkg -case=underscore -testonly

// ErrorPresenter ...
type ErrorPresenter interface {
	PresentError(
		writer http.ResponseWriter,
		request *http.Request,
		statusCode int,
		err error,
	) error
}

// SilentErrorPresenter ...
type SilentErrorPresenter struct {
	ErrorPresenter ErrorPresenter
	Logger         log.Logger
}

// PresentError ...
func (presenter SilentErrorPresenter) PresentError(
	writer http.ResponseWriter,
	request *http.Request,
	statusCode int,
	err error,
) {
	err = presenter.ErrorPresenter.PresentError(writer, request, statusCode, err)
	if err != nil {
		presenter.Logger.Logf("unable to present the error: %v", err)
	}
}
