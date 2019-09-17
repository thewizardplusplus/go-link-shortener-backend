package usecases

import (
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

// LinkGetter ...
type LinkGetter interface {
	GetLink(query string) (entities.Link, error)
}
