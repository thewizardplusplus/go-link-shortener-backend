package httputils

import (
	"net/http"

	"github.com/gorilla/mux"
)

//go:generate mockery -name=Printer -inpkg -case=underscore -testonly

// Printer ...
type Printer interface {
	Printf(template string, arguments ...interface{})
}

// CatchingMiddleware ...
func CatchingMiddleware(printer Printer) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			catchingWriter := NewCatchingResponseWriter(writer)
			next.ServeHTTP(catchingWriter, request)

			if err := catchingWriter.LastError(); err != nil {
				printer.Printf("unable to handle the request: %v", err)
			}
		})
	}
}
