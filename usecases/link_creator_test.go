package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

func TestLinkCreator_CreateLink(test *testing.T) {
	type fields struct {
		LinkGetter    LinkGetter
		LinkSetter    LinkSetter
		CodeGenerator CodeGenerator
	}
	type args struct {
		url string
	}

	for _, data := range []struct {
		name     string
		fields   fields
		args     args
		wantLink entities.Link
		wantErr  assert.ErrorAssertionFunc
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			creator := LinkCreator{
				LinkGetter:    data.fields.LinkGetter,
				LinkSetter:    data.fields.LinkSetter,
				CodeGenerator: data.fields.CodeGenerator,
			}
			gotLink, gotErr := creator.CreateLink(data.args.url)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.LinkGetter,
				data.fields.LinkSetter,
				data.fields.CodeGenerator,
			)
			assert.Equal(test, data.wantLink, gotLink)
			data.wantErr(test, gotErr)
		})
	}
}
