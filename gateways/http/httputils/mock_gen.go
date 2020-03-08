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
