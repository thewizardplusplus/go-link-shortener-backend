package presenters

import (
	"encoding/json"
	"io"
	"net/http"
)

// JSONPresenter ...
type JSONPresenter struct{}

func presentData(writer http.ResponseWriter, statusCode int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	bytes, _ := json.Marshal(data)        // nolint: gosec
	io.WriteString(writer, string(bytes)) // nolint: gosec, errcheck
}
