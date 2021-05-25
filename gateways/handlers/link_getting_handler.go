package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
)

//go:generate mockery --name=LinkGetter --inpackage --case=underscore --testonly

// LinkGetter ...
type LinkGetter interface {
	GetLink(code string) (entities.Link, error)
}

//go:generate mockery --name=LinkPresenter --inpackage --case=underscore --testonly

// LinkPresenter ...
type LinkPresenter interface {
	PresentLink(
		writer http.ResponseWriter,
		request *http.Request,
		link entities.Link,
	)
}

//go:generate mockery --name=ErrorPresenter --inpackage --case=underscore --testonly

// ErrorPresenter ...
type ErrorPresenter interface {
	PresentError(
		writer http.ResponseWriter,
		request *http.Request,
		statusCode int,
		err error,
	)
}

// LinkGettingHandler ...
type LinkGettingHandler struct {
	LinkGetter     LinkGetter
	LinkPresenter  LinkPresenter
	ErrorPresenter ErrorPresenter
}

// ServeHTTP ...
//   @router /links/{code} [GET]
//   @param code path string true "link code"
//   @produce json
//   @success 200 {object} entities.Link
//   @failure 404 {object} presenters.ErrorResponse
//   @failure 500 {object} presenters.ErrorResponse
func (handler LinkGettingHandler) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
) {
	code := mux.Vars(request)["code"]
	link, err := handler.LinkGetter.GetLink(code)
	switch err {
	case nil:
		handler.LinkPresenter.PresentLink(writer, request, link)
	case sql.ErrNoRows:
		const statusCode = http.StatusNotFound
		err = errors.New("unable to find the link")
		handler.ErrorPresenter.PresentError(writer, request, statusCode, err)
	default:
		const statusCode = http.StatusInternalServerError
		err = errors.Wrap(err, "unable to get the link")
		handler.ErrorPresenter.PresentError(writer, request, statusCode, err)
	}
}
