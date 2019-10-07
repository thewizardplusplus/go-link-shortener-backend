package handlers

import (
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

// LinkCreator ...
type LinkCreator interface {
	CreateLink(url string) (entities.Link, error)
}
