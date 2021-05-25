package handlers

import (
	"net/http"
)

//go:generate mockery --name=Handler --inpackage --case=underscore --testonly

// Handler ...
//
// It is used only for mock generating.
//
type Handler interface {
	http.Handler
}
