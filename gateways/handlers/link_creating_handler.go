package handlers

import (
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

//go:generate mockery -name=LinkCreator -inpkg -case=underscore -testonly

// LinkCreator ...
type LinkCreator interface {
	CreateLink(url string) (entities.Link, error)
}
