package presenters

import (
	"net/http"

	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
)

//go:generate mockery --name=LinkPresenter --inpackage --case=underscore --testonly

// LinkPresenter ...
type LinkPresenter interface {
	PresentLink(
		writer http.ResponseWriter,
		request *http.Request,
		link entities.Link,
	) error
}

// SilentLinkPresenter ...
type SilentLinkPresenter struct {
	LinkPresenter LinkPresenter
	Logger        log.Logger
}

// PresentLink ...
func (presenter SilentLinkPresenter) PresentLink(
	writer http.ResponseWriter,
	request *http.Request,
	link entities.Link,
) {
	err := presenter.LinkPresenter.PresentLink(writer, request, link)
	if err != nil {
		presenter.Logger.Logf("unable to present the link: %v", err)
	}
}
