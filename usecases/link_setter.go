package usecases

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
)

//go:generate mockery -name=LinkSetter -inpkg -case=underscore -testonly

// LinkSetter ...
type LinkSetter interface {
	SetLink(link entities.Link) error
}

// SilentLinkSetter ...
type SilentLinkSetter struct {
	LinkSetter LinkSetter
	Printer    Printer
}

// SetLink ...
func (setter SilentLinkSetter) SetLink(link entities.Link) error {
	if err := setter.LinkSetter.SetLink(link); err != nil {
		setter.Printer.Printf("unable to set the link: %v", err)
	}

	return nil
}

// LinkSetterGroup ...
type LinkSetterGroup []LinkSetter

// SetLink ...
func (setters LinkSetterGroup) SetLink(link entities.Link) error {
	for _, setter := range setters {
		if err := setter.SetLink(link); err != nil {
			return errors.Wrap(err, "unable to set the link")
		}
	}

	return nil
}
