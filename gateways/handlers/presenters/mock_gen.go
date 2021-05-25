package presenters

import (
	"net/http"

	"github.com/go-log/log"
)

//go:generate mockery --name=Logger --inpackage --case=underscore --testonly

// Logger ...
//
// It is used only for mock generating.
//
type Logger interface {
	log.Logger
}

//go:generate mockery --name=ResponseWriter --inpackage --case=underscore --testonly

// ResponseWriter ...
//
// It is used only for mock generating.
//
type ResponseWriter interface {
	http.ResponseWriter
}
