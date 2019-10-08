package handlers

import (
	"encoding/json"
	"net/http"

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

type linkCreatingRequest struct {
	URL string
}

// ServeHTTP ...
func (handler LinkCreatingHandler) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var data linkCreatingRequest
	if err := json.NewDecoder(request.Body).Decode(&data); err != nil {
		const statusCode = http.StatusBadRequest
		err = errors.Wrap(err, "unable to decode the request")
		handler.ErrorPresenter.PresentError(writer, statusCode, err)

		return
	}

	link, err := handler.LinkCreator.CreateLink(data.URL)
	if err != nil {
		const statusCode = http.StatusInternalServerError
		err = errors.Wrap(err, "unable to create the link")
		handler.ErrorPresenter.PresentError(writer, statusCode, err)

		return
	}

	handler.LinkPresenter.PresentLink(writer, link)
}
