package httputils

import (
	"net/http"
)

//go:generate mockery -name=Printer -inpkg -case=underscore -testonly

// Printer ...
type Printer interface {
	Printf(template string, arguments ...interface{})
}

// CatchingHandler ...
type CatchingHandler struct {
	Handler http.Handler
	Printer Printer
}

// ServeHTTP ...
func (handler CatchingHandler) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
) {
	catchingWriter := NewCatchingResponseWriter(writer)
	handler.Handler.ServeHTTP(catchingWriter, request)

	if err := catchingWriter.LastError(); err != nil {
		handler.Printer.Printf("unable to handle the request: %v", err)
	}
}
