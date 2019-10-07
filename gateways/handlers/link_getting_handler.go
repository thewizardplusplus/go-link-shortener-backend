package handlers

import (
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

//go:generate mockery -name=LinkGetter -inpkg -case=underscore -testonly

// LinkGetter ...
type LinkGetter interface {
	GetLink(code string) (entities.Link, error)
}
