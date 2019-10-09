package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handlers ...
type Handlers struct {
	LinkGettingHandler  http.Handler
	LinkCreatingHandler http.Handler
	NotFoundHandler     http.Handler
}

// NewRouter ...
func NewRouter(handlers Handlers) http.Handler {
	router := mux.NewRouter().PathPrefix("/api/v1/links").Subrouter()
	router.NotFoundHandler = handlers.NotFoundHandler
	router.MethodNotAllowedHandler = handlers.NotFoundHandler
	router.Handle("/{code}", handlers.LinkGettingHandler).Methods(http.MethodGet)
	router.Handle("/", handlers.LinkCreatingHandler).Methods(http.MethodPost)

	return router
}
