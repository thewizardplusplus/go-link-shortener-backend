package presenters

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

func TestSilentLinkPresenter_PresentLink(test *testing.T) {
	type fields struct {
		LinkPresenter LinkPresenter
		Printer       Printer
	}
	type args struct {
		writer http.ResponseWriter
		link   entities.Link
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(t *testing.T) {
			presenter := SilentLinkPresenter{
				LinkPresenter: data.fields.LinkPresenter,
				Printer:       data.fields.Printer,
			}
			presenter.PresentLink(data.args.writer, data.args.link)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.LinkPresenter,
				data.fields.Printer,
				data.args.writer,
			)
		})
	}
}
