package usecases

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
)

//go:generate mockery -name=CodeGenerator -inpkg -case=underscore -testonly

// CodeGenerator ...
type CodeGenerator interface {
	GenerateCode() (string, error)
}

// LinkCreator ...
type LinkCreator struct {
	LinkGetter    LinkGetter
	LinkSetter    LinkSetter
	CodeGenerator CodeGenerator
}

// CreateLink ...
func (creator LinkCreator) CreateLink(url string) (entities.Link, error) {
	link, err := creator.LinkGetter.GetLink(url)
	switch err {
	case nil:
		return link, nil
	case sql.ErrNoRows:
	default:
		return entities.Link{}, errors.Wrap(err, "unable to get the link")
	}

	code, err := creator.CodeGenerator.GenerateCode()
	if err != nil {
		return entities.Link{}, errors.Wrap(err, "unable to generate a code")
	}

	link = entities.Link{Code: code, URL: url}
	if err := creator.LinkSetter.SetLink(link); err != nil {
		return entities.Link{}, errors.Wrap(err, "unable to set the link")
	}

	return link, nil
}
