package httputils

import (
	"net/http"
)

// CatchingRedirect ...
func CatchingRedirect(
	writer http.ResponseWriter,
	request *http.Request,
	url string,
	statusCode int,
) error {
	catchingWriter := NewCatchingResponseWriter(writer)
	http.Redirect(catchingWriter, request, url, statusCode)

	return catchingWriter.LastError()
}
