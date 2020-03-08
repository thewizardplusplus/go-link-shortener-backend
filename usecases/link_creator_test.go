package usecases

import (
	"database/sql"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
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
		{
			name: "success with the getter",
			fields: fields{
				LinkGetter: func() LinkGetter {
					getter := new(MockLinkGetter)
					getter.
						On("GetLink", "url").
						Return(entities.Link{Code: "code", URL: "url"}, nil)

					return getter
				}(),
				LinkSetter:    new(MockLinkSetter),
				CodeGenerator: new(MockCodeGenerator),
			},
			args:     args{"url"},
			wantLink: entities.Link{Code: "code", URL: "url"},
			wantErr:  assert.NoError,
		},
		{
			name: "success with the setter",
			fields: fields{
				LinkGetter: func() LinkGetter {
					getter := new(MockLinkGetter)
					getter.On("GetLink", "url").Return(entities.Link{}, sql.ErrNoRows)

					return getter
				}(),
				LinkSetter: func() LinkSetter {
					setter := new(MockLinkSetter)
					setter.On("SetLink", entities.Link{Code: "code", URL: "url"}).Return(nil)

					return setter
				}(),
				CodeGenerator: func() CodeGenerator {
					generator := new(MockCodeGenerator)
					generator.On("GenerateCode").Return("code", nil)

					return generator
				}(),
			},
			args:     args{"url"},
			wantLink: entities.Link{Code: "code", URL: "url"},
			wantErr:  assert.NoError,
		},
		{
			name: "error with the getter",
			fields: fields{
				LinkGetter: func() LinkGetter {
					getter := new(MockLinkGetter)
					getter.On("GetLink", "url").Return(entities.Link{}, iotest.ErrTimeout)

					return getter
				}(),
				LinkSetter:    new(MockLinkSetter),
				CodeGenerator: new(MockCodeGenerator),
			},
			args:     args{"url"},
			wantLink: entities.Link{},
			wantErr:  assert.Error,
		},
		{
			name: "error with the generator",
			fields: fields{
				LinkGetter: func() LinkGetter {
					getter := new(MockLinkGetter)
					getter.On("GetLink", "url").Return(entities.Link{}, sql.ErrNoRows)

					return getter
				}(),
				LinkSetter: new(MockLinkSetter),
				CodeGenerator: func() CodeGenerator {
					generator := new(MockCodeGenerator)
					generator.On("GenerateCode").Return("", iotest.ErrTimeout)

					return generator
				}(),
			},
			args:     args{"url"},
			wantLink: entities.Link{},
			wantErr:  assert.Error,
		},
		{
			name: "error with the setter",
			fields: fields{
				LinkGetter: func() LinkGetter {
					getter := new(MockLinkGetter)
					getter.On("GetLink", "url").Return(entities.Link{}, sql.ErrNoRows)

					return getter
				}(),
				LinkSetter: func() LinkSetter {
					setter := new(MockLinkSetter)
					setter.
						On("SetLink", entities.Link{Code: "code", URL: "url"}).
						Return(iotest.ErrTimeout)

					return setter
				}(),
				CodeGenerator: func() CodeGenerator {
					generator := new(MockCodeGenerator)
					generator.On("GenerateCode").Return("code", nil)

					return generator
				}(),
			},
			args:     args{"url"},
			wantLink: entities.Link{},
			wantErr:  assert.Error,
		},
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
