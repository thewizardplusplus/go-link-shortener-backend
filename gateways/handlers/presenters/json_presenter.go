package presenters

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
	httputils "github.com/thewizardplusplus/go-http-utils"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
)

// JSONPresenter ...
type JSONPresenter struct{}

// ErrorResponse ...
//
// It's public only for docs generating.
type ErrorResponse struct {
	Error string
}

// PresentLink ...
func (presenter JSONPresenter) PresentLink(
	writer http.ResponseWriter,
	request *http.Request,
	link entities.Link,
) error {
	if err := httputils.WriteJSON(writer, http.StatusOK, link); err != nil {
		return errors.Wrap(err, "unable to present the link in JSON")
	}

	return nil
}

// PresentError ...
func (presenter JSONPresenter) PresentError(
	writer http.ResponseWriter,
	request *http.Request,
	statusCode int,
	err error,
) error {
	response := ErrorResponse{Error: err.Error()}
	if err := httputils.WriteJSON(writer, statusCode, response); err != nil {
		return errors.Wrap(err, "unable to present the error in JSON")
	}

	return nil
}

func presentData(
	writer http.ResponseWriter,
	statusCode int,
	data interface{},
) error {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	bytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "unable to marshal the data")
	}

	if _, err := io.WriteString(writer, string(bytes)); err != nil {
		return errors.Wrap(err, "unable to write the data")
	}

	return nil
}
