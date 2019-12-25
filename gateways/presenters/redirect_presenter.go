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
