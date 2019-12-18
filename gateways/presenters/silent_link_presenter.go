package presenters

import (
	"net/http"

	"github.com/thewizardplusplus/go-link-shortener/entities"
)

// LinkPresenter ...
type LinkPresenter interface {
	PresentLink(writer http.ResponseWriter, link entities.Link) error
}
