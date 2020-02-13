package presenters

import (
	"net/http"

	"github.com/pkg/errors"
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
) error {
	err := redirect(writer, request, link.URL, http.StatusMovedPermanently)
	if err != nil {
		return errors.Wrap(err, "unable to redirect to the link")
	}

	return nil
}

// PresentError ...
func (presenter RedirectPresenter) PresentError(
	writer http.ResponseWriter,
	request *http.Request,
	statusCode int,
	err error,
) error {
	err2 := redirect(writer, request, presenter.ErrorURL, http.StatusFound)
	if err2 != nil {
		return errors.Wrap(err2, "unable to redirect to the error")
	}

	presenter.Printer.Printf("redirect because of the error: %v", err)
	return nil
}

func redirect(
	writer http.ResponseWriter,
	request *http.Request,
	url string,
	statusCode int,
) error {
	catchingWriter := newCatchingResponseWriter(writer)
	http.Redirect(catchingWriter, request, url, statusCode)

	// errors with writing to the http.ResponseWriter is important to handle,
	// see for details: https://stackoverflow.com/a/43976633
	if catchingWriter.error != nil {
		return errors.Wrap(catchingWriter.error, "unable to write the data")
	}

	return nil
}
