package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

//go:generate mockery -name=LinkCreator -inpkg -case=underscore -testonly

// LinkCreator ...
type LinkCreator interface {
	CreateLink(url string) (entities.Link, error)
}

// LinkCreatingHandler ...
type LinkCreatingHandler struct {
	LinkCreator    LinkCreator
	LinkPresenter  LinkPresenter
	ErrorPresenter ErrorPresenter
}

// ServeHTTP ...
func (handler LinkCreatingHandler) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
) {
	url := mux.Vars(request)["url"]
	link, err := handler.LinkCreator.CreateLink(url)
	if err != nil {
		const statusCode = http.StatusInternalServerError
		err = errors.Wrap(err, "unable to create the link")
		handler.ErrorPresenter.PresentError(writer, statusCode, err)

		return
	}

	handler.LinkPresenter.PresentLink(writer, link)
}
