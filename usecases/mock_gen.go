package usecases

import (
	"github.com/go-log/log"
)

//go:generate mockery --name=Logger --inpackage --case=underscore --testonly

// Logger ...
//
// It's used only for mock generating.
type Logger interface {
	log.Logger
}
