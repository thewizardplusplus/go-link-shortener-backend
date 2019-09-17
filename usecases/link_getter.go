package usecases

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

//go:generate mockery -name=LinkGetter -inpkg -case=underscore -testonly

// LinkGetter ...
type LinkGetter interface {
	GetLink(query string) (entities.Link, error)
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
