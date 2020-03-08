package presenters

import (
	"net/http"

	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
)

//go:generate mockery -name=LinkPresenter -inpkg -case=underscore -testonly

// LinkPresenter ...
type LinkPresenter interface {
	PresentLink(
		writer http.ResponseWriter,
		request *http.Request,
		link entities.Link,
	) error
}

//go:generate mockery -name=Printer -inpkg -case=underscore -testonly

// Printer ...
type Printer interface {
	Printf(template string, arguments ...interface{})
}

// SilentLinkPresenter ...
type SilentLinkPresenter struct {
	LinkPresenter LinkPresenter
	Printer       Printer
}

// PresentLink ...
func (presenter SilentLinkPresenter) PresentLink(
	writer http.ResponseWriter,
	request *http.Request,
	link entities.Link,
) {
	err := presenter.LinkPresenter.PresentLink(writer, request, link)
	if err != nil {
		presenter.Printer.Printf("unable to present the link: %v", err)
	}
}
