package httputils

import (
	"net/http"

	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
)

// SPAFallbackMiddleware ...
func SPAFallbackMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			if isStaticAssetRequest(request) {
				request.URL.Path = "/"
			}
			next.ServeHTTP(writer, request)
		})
	}
}

func isStaticAssetRequest(request *http.Request) bool {
	if request.Method != http.MethodGet {
		return false
	}

	for _, spec := range header.ParseAccept(request.Header, "Accept") {
		if spec.Value == "text/html" {
			return true
		}
	}

	return false
}
