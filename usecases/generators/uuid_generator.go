package generators

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// UUIDGenerator ...
type UUIDGenerator struct{}

// GenerateCode ...
func (generator UUIDGenerator) GenerateCode() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "unable to generate an UUID V4")
	}

	return uuid.String(), nil
}
