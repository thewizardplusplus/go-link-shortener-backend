package usecases

import (
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

// LinkSetter ...
type LinkSetter interface {
	SetLink(link entities.Link) error
}
