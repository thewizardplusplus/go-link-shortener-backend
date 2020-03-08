package router

import (
	"net/http"
)

//go:generate mockery -name=Handler -inpkg -case=underscore -testonly

// Handler ...
//
// It's used only for mock generating.
type Handler interface {
	http.Handler
}
