package presenters

import (
	"net/http"

	"github.com/thewizardplusplus/go-link-shortener/entities"
)

// RedirectPresenter ...
type RedirectPresenter struct {
	ErrorURL string
	Printer  Printer
}

// PresentLink ...
func (presenter RedirectPresenter) PresentLink(
	writer http.ResponseWriter,
	request *http.Request,
	link entities.Link,
) {
	http.Redirect(writer, request, link.URL, http.StatusMovedPermanently)
}

// PresentError ...
func (presenter RedirectPresenter) PresentError(
	writer http.ResponseWriter,
	request *http.Request,
	statusCode int,
	err error,
) {
	http.Redirect(writer, request, presenter.ErrorURL, http.StatusFound)
	presenter.Printer.Printf("redirect because of the error: %v", err)
}
