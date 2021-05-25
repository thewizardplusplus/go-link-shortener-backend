package handlers

import (
	"net/http"

	"github.com/pkg/errors"
	httputils "github.com/thewizardplusplus/go-http-utils"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
)

//go:generate mockery --name=LinkCreator --inpackage --case=underscore --testonly

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

// LinkCreatingRequest ...
//
// It's public only for docs generating.
type LinkCreatingRequest struct {
	URL string
}

// ServeHTTP ...
//   @router /links/ [POST]
//   @accept json
//   @param data body handlers.LinkCreatingRequest true "link data"
//   @produce json
//   @success 200 {object} entities.Link
//   @failure 400 {object} presenters.ErrorResponse
//   @failure 500 {object} presenters.ErrorResponse
func (handler LinkCreatingHandler) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var data LinkCreatingRequest
	if err := httputils.ReadJSON(request.Body, &data); err != nil {
		const statusCode = http.StatusBadRequest
		err = errors.Wrap(err, "unable to decode the request")
		handler.ErrorPresenter.PresentError(writer, request, statusCode, err)

		return
	}

	link, err := handler.LinkCreator.CreateLink(data.URL)
	if err != nil {
		const statusCode = http.StatusInternalServerError
		err = errors.Wrap(err, "unable to create the link")
		handler.ErrorPresenter.PresentError(writer, request, statusCode, err)

		return
	}

	handler.LinkPresenter.PresentLink(writer, request, link)
}
