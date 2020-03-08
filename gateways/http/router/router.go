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
	StaticFileHandler   http.Handler
}

// NewRouter ...
func NewRouter(redirectEndpointPrefix string, handlers Handlers) *mux.Router {
	// @title go-link-shortener API
	// @version 1.0.0
	// @license.name MIT
	// @host localhost:8080
	// @basePath /api/v1

	rootRouter := mux.NewRouter()

	apiRouter := rootRouter.PathPrefix("/api").Subrouter()
	apiV1Router := apiRouter.PathPrefix("/v1").Subrouter()
	linksAPIV1Router := apiV1Router.PathPrefix("/links").Subrouter()
	linksAPIV1Router.
		Handle("/{code}", handlers.LinkGettingHandler).
		Methods(http.MethodGet)
	linksAPIV1Router.
		Handle("/", handlers.LinkCreatingHandler).
		Methods(http.MethodPost)

	rootRouter.
		Handle(redirectEndpointPrefix+"/{code}", handlers.LinkRedirectHandler)
	rootRouter.
		PathPrefix("/").Handler(handlers.StaticFileHandler).
		Methods(http.MethodGet)

	return rootRouter
}
