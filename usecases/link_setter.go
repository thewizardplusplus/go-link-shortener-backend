package usecases

import (
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

//go:generate mockery -name=LinkSetter -inpkg -case=underscore -testonly

// LinkSetter ...
type LinkSetter interface {
	SetLink(link entities.Link) error
}

// LinkSetterGroup ...
type LinkSetterGroup []LinkSetter
