package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

func TestLinkSetterGroup_SetLink(test *testing.T) {
	type args struct {
		link entities.Link
	}

	for _, data := range []struct {
		name    string
		setters LinkSetterGroup
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			gotErr := data.setters.SetLink(data.args.link)

			for _, setter := range data.setters {
				mock.AssertExpectationsForObjects(test, setter)
			}
			data.wantErr(test, gotErr)
		})
	}
}
