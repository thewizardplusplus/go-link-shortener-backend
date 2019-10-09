package presenters

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/thewizardplusplus/go-link-shortener/entities"
)

// JSONPresenter ...
type JSONPresenter struct{}

// PresentLink ...
func (presenter JSONPresenter) PresentLink(
	writer http.ResponseWriter,
	link entities.Link,
) {
	presentData(writer, http.StatusOK, link)
}

func presentData(writer http.ResponseWriter, statusCode int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	bytes, _ := json.Marshal(data)        // nolint: gosec
	io.WriteString(writer, string(bytes)) // nolint: gosec, errcheck
}
