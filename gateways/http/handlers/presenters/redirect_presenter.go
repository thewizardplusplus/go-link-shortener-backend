package presenters

// nolint: lll
import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
	"github.com/thewizardplusplus/go-link-shortener-backend/gateways/http/httputils"
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
	url, statusCode := link.URL, http.StatusMovedPermanently
	err := httputils.CatchingRedirect(writer, request, url, statusCode)
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
	url, statusCode2 := presenter.ErrorURL, http.StatusFound
	err2 := httputils.CatchingRedirect(writer, request, url, statusCode2)
	if err2 != nil {
		return errors.Wrap(err2, "unable to redirect to the error")
	}

	presenter.Printer.Printf("redirect because of the error: %v", err)
	return nil
}
