package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handlers ...
type Handlers struct {
	LinkRedirectHandler http.Handler
	LinkGettingHandler  http.Handler
	LinkCreatingHandler http.Handler
	StaticFileHandler   http.Handler
}

// NewRouter ...
func NewRouter(redirectEndpointPrefix string, handlers Handlers) *mux.Router {
	// @title go-link-shortener API
	// @version 1.9.0
	// @license.name MIT
	// @host localhost:8080
	// @basePath /api/v1

	rootRouter := mux.NewRouter()
	apiRouter := rootRouter.PathPrefix("/api/v1").Subrouter()

	rootRouter.
		Handle(redirectEndpointPrefix+"/{code}", handlers.LinkRedirectHandler)
	rootRouter.
		PathPrefix("/").Handler(handlers.StaticFileHandler).
		Methods(http.MethodGet)

	apiRouter.
		Handle("/links/{code}", handlers.LinkGettingHandler).
		Methods(http.MethodGet)
	apiRouter.
		Handle("/links/", handlers.LinkCreatingHandler).
		Methods(http.MethodPost)

	return rootRouter
}
