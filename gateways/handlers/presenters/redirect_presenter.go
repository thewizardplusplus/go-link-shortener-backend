package presenters

// nolint: lll
import (
	"net/http"

	"github.com/go-log/log"
	"github.com/pkg/errors"
	httputils "github.com/thewizardplusplus/go-http-utils"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
)

// RedirectPresenter ...
type RedirectPresenter struct {
	ErrorURL string
	Logger   log.Logger
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

	presenter.Logger.Logf("redirect because of the error: %v", err)
	return nil
}
