package httputils

import (
	"net/http"
)

//go:generate mockery -name=ResponseWriter -inpkg -case=underscore -testonly

// ResponseWriter ...
//
// It's used only for mock generating.
type ResponseWriter interface {
	http.ResponseWriter
}

//go:generate mockery -name=Handler -inpkg -case=underscore -testonly

// Handler ...
//
// It's used only for mock generating.
type Handler interface {
	http.Handler
}
