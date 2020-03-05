package httputils

import (
	"net/http"
)

// CatchingResponseWriter ...
type CatchingResponseWriter struct {
	http.ResponseWriter

	lastError error
}

// NewCatchingResponseWriter ...
func NewCatchingResponseWriter(
	writer http.ResponseWriter,
) *CatchingResponseWriter {
	return &CatchingResponseWriter{ResponseWriter: writer}
}

// LastError ...
func (writer CatchingResponseWriter) LastError() error {
	return writer.lastError
}

// Write ...
func (writer *CatchingResponseWriter) Write(p []byte) (n int, err error) {
	n, err = writer.ResponseWriter.Write(p)
	writer.lastError = err

	return n, err
}
