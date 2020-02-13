package presenters

import (
	"net/http"
)

type catchingResponseWriter struct {
	http.ResponseWriter

	error error
}

func newCatchingResponseWriter(
	writer http.ResponseWriter,
) *catchingResponseWriter {
	return &catchingResponseWriter{ResponseWriter: writer}
}

func (writer *catchingResponseWriter) Write(p []byte) (n int, err error) {
	n, err = writer.ResponseWriter.Write(p)
	writer.error = err

	return n, err
}
