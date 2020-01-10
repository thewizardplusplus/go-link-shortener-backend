package presenters

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

func TestRedirectPresenter_PresentLink(test *testing.T) {
	type fields struct {
		ErrorURL string
		Printer  Printer
	}
	type args struct {
		writer http.ResponseWriter
		link   entities.Link
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
		check   func(test *testing.T, writer http.ResponseWriter)
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			presenter := RedirectPresenter{
				ErrorURL: data.fields.ErrorURL,
				Printer:  data.fields.Printer,
			}
			gotErr := presenter.PresentLink(data.args.writer, data.args.link)

			mock.AssertExpectationsForObjects(test, data.fields.Printer)
			data.wantErr(test, gotErr)
			data.check(test, data.args.writer)
		})
	}
}
