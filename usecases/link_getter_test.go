package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

func TestLinkGetterGroup_GetLink(test *testing.T) {
	type args struct {
		query string
	}

	for _, data := range []struct {
		name     string
		getters  LinkGetterGroup
		args     args
		wantLink entities.Link
		wantErr  assert.ErrorAssertionFunc
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			gotLink, gotErr := data.getters.GetLink(data.args.query)

			for _, getter := range data.getters {
				mock.AssertExpectationsForObjects(test, getter)
			}
			assert.Equal(test, data.wantLink, gotLink)
			data.wantErr(test, gotErr)
		})
	}
}
