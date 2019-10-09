package router

import (
	"net/http"
)

// Handlers ...
type Handlers struct {
	LinkGettingHandler  http.Handler
	LinkCreatingHandler http.Handler
	NotFoundHandler     http.Handler
}
