package presenters

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-link-shortener/entities"
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
	link entities.Link,
) {
	presentData(writer, http.StatusOK, link)
}

// PresentError ...
func (presenter JSONPresenter) PresentError(
	writer http.ResponseWriter,
	statusCode int,
	err error,
) {
	response := ErrorResponse{Error: err.Error()}
	presentData(writer, statusCode, response)
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
