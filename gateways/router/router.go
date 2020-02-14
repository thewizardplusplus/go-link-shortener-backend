package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handlers ...
type Handlers struct {
	LinkRedirectHandler http.Handler
	LinkGettingHandler  http.Handler
	LinkCreatingHandler http.Handler
	NotFoundHandler     http.Handler
}

// NewRouter ...
func NewRouter(redirectEndpointPrefix string, handlers Handlers) *mux.Router {
	// @title go-link-shortener API
	// @version 1.0.0
	// @license.name MIT
	// @host localhost:8080
	// @basePath /api/v1

	router := mux.NewRouter()
	router.NotFoundHandler = handlers.NotFoundHandler
	router.MethodNotAllowedHandler = handlers.NotFoundHandler
	router.Handle(redirectEndpointPrefix+"/{code}", handlers.LinkRedirectHandler)

	apiRouter := router.PathPrefix("/api/v1/links").Subrouter()
	apiRouter.
		Handle("/{code}", handlers.LinkGettingHandler).
		Methods(http.MethodGet)
	apiRouter.Handle("/", handlers.LinkCreatingHandler).Methods(http.MethodPost)

	return router
}
