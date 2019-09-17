package usecases

import (
	"testing"
	"testing/iotest"

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
		{
			name:    "success without setters",
			setters: nil,
			args: args{
				link: entities.Link{Code: "code", URL: "url"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with setters",
			setters: func() LinkSetterGroup {
				setterOne := new(MockLinkSetter)
				setterOne.On("SetLink", entities.Link{Code: "code", URL: "url"}).Return(nil)

				setterTwo := new(MockLinkSetter)
				setterTwo.On("SetLink", entities.Link{Code: "code", URL: "url"}).Return(nil)

				return LinkSetterGroup{setterOne, setterTwo}
			}(),
			args: args{
				link: entities.Link{Code: "code", URL: "url"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error with the first setter",
			setters: func() LinkSetterGroup {
				setterOne := new(MockLinkSetter)
				setterOne.
					On("SetLink", entities.Link{Code: "code", URL: "url"}).
					Return(iotest.ErrTimeout)

				setterTwo := new(MockLinkSetter)

				return LinkSetterGroup{setterOne, setterTwo}
			}(),
			args: args{
				link: entities.Link{Code: "code", URL: "url"},
			},
			wantErr: assert.Error,
		},
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
