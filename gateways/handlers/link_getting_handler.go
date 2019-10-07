package handlers

import (
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

// LinkGetter ...
type LinkGetter interface {
	GetLink(code string) (entities.Link, error)
}
