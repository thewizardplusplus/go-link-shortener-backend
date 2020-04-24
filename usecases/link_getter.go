package usecases

import (
	"database/sql"

	"github.com/go-log/log"
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
)

//go:generate mockery -name=LinkGetter -inpkg -case=underscore -testonly

// LinkGetter ...
type LinkGetter interface {
	GetLink(query string) (entities.Link, error)
}

// SilentLinkGetter ...
type SilentLinkGetter struct {
	LinkGetter LinkGetter
	Logger     log.Logger
}

// GetLink ...
func (getter SilentLinkGetter) GetLink(query string) (entities.Link, error) {
	link, err := getter.LinkGetter.GetLink(query)
	switch err {
	case nil:
		return link, nil
	default:
		if err != sql.ErrNoRows {
			getter.Logger.Logf("unable to get the link: %v", err)
		}

		return entities.Link{}, sql.ErrNoRows
	}
}

// LinkGetterGroup ...
type LinkGetterGroup []LinkGetter

// GetLink ...
func (getters LinkGetterGroup) GetLink(query string) (entities.Link, error) {
	for _, getter := range getters {
		link, err := getter.GetLink(query)
		switch err {
		case nil:
			return link, nil
		case sql.ErrNoRows:
		default:
			return entities.Link{}, errors.Wrap(err, "unable to get the link")
		}
	}

	return entities.Link{}, sql.ErrNoRows
}
