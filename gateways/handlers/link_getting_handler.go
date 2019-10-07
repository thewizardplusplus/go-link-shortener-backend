package handlers

import (
	"net/http"

	"github.com/thewizardplusplus/go-link-shortener/entities"
)

//go:generate mockery -name=LinkGetter -inpkg -case=underscore -testonly

// LinkGetter ...
type LinkGetter interface {
	GetLink(code string) (entities.Link, error)
}

//go:generate mockery -name=LinkPresenter -inpkg -case=underscore -testonly

// LinkPresenter ...
type LinkPresenter interface {
	PresentLink(writer http.ResponseWriter, link entities.Link)
}
