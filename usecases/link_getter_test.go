package usecases

import (
	"database/sql"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

func TestSilentLinkGetter_GetLink(test *testing.T) {
	type fields struct {
		LinkGetter LinkGetter
		Printer    Printer
	}
	type args struct {
		query string
	}

	for _, data := range []struct {
		name     string
		fields   fields
		args     args
		wantLink entities.Link
		wantErr  error
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			getter := SilentLinkGetter{
				LinkGetter: data.fields.LinkGetter,
				Printer:    data.fields.Printer,
			}
			gotLink, gotErr := getter.GetLink(data.args.query)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.LinkGetter,
				data.fields.Printer,
			)
			assert.Equal(test, data.wantLink, gotLink)
			assert.Equal(test, data.wantErr, gotErr)
		})
	}
}

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
		{
			name: "success with the first getter",
			getters: func() LinkGetterGroup {
				getterOne := new(MockLinkGetter)
				getterOne.
					On("GetLink", "query").
					Return(entities.Link{Code: "code", URL: "url"}, nil)

				getterTwo := new(MockLinkGetter)

				return LinkGetterGroup{getterOne, getterTwo}
			}(),
			args:     args{"query"},
			wantLink: entities.Link{Code: "code", URL: "url"},
			wantErr:  assert.NoError,
		},
		{
			name: "success not with the first getter",
			getters: func() LinkGetterGroup {
				getterOne := new(MockLinkGetter)
				getterOne.On("GetLink", "query").Return(entities.Link{}, sql.ErrNoRows)

				getterTwo := new(MockLinkGetter)
				getterTwo.
					On("GetLink", "query").
					Return(entities.Link{Code: "code", URL: "url"}, nil)

				return LinkGetterGroup{getterOne, getterTwo}
			}(),
			args:     args{"query"},
			wantLink: entities.Link{Code: "code", URL: "url"},
			wantErr:  assert.NoError,
		},
		{
			name:     "error without getters",
			getters:  nil,
			args:     args{"query"},
			wantLink: entities.Link{},
			wantErr: func(test assert.TestingT, err error, args ...interface{}) bool {
				return assert.Equal(test, sql.ErrNoRows, err, args)
			},
		},
		{
			name: "error with the first getter",
			getters: func() LinkGetterGroup {
				getterOne := new(MockLinkGetter)
				getterOne.On("GetLink", "query").Return(entities.Link{}, iotest.ErrTimeout)

				getterTwo := new(MockLinkGetter)

				return LinkGetterGroup{getterOne, getterTwo}
			}(),
			args:     args{"query"},
			wantLink: entities.Link{},
			wantErr:  assert.Error,
		},
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
