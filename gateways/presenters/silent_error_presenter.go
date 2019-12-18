package presenters

import (
	"net/http"
)

//go:generate mockery -name=ErrorPresenter -inpkg -case=underscore -testonly

// ErrorPresenter ...
type ErrorPresenter interface {
	PresentError(writer http.ResponseWriter, statusCode int, err error) error
}

// SilentErrorPresenter ...
type SilentErrorPresenter struct {
	ErrorPresenter ErrorPresenter
	Printer        Printer
}

// PresentError ...
func (presenter SilentErrorPresenter) PresentError(
	writer http.ResponseWriter,
	statusCode int,
	err error,
) {
	if err = presenter.ErrorPresenter.PresentError(
		writer,
		statusCode,
		err,
	); err != nil {
		presenter.Printer.Printf("unable to present the error: %v", err)
	}
}
