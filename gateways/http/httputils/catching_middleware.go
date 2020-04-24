package httputils

import (
	"net/http"

	"github.com/go-log/log"
	"github.com/gorilla/mux"
)

// CatchingMiddleware ...
func CatchingMiddleware(logger log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			catchingWriter := NewCatchingResponseWriter(writer)
			next.ServeHTTP(catchingWriter, request)

			if err := catchingWriter.LastError(); err != nil {
				logger.Logf("unable to handle the request: %v", err)
			}
		})
	}
}
