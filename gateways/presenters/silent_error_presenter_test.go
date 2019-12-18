package presenters

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestSilentErrorPresenter_PresentError(test *testing.T) {
	type fields struct {
		ErrorPresenter ErrorPresenter
		Printer        Printer
	}
	type args struct {
		writer     http.ResponseWriter
		statusCode int
		err        error
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			presenter := SilentErrorPresenter{
				ErrorPresenter: data.fields.ErrorPresenter,
				Printer:        data.fields.Printer,
			}
			presenter.PresentError(data.args.writer, data.args.statusCode, data.args.err)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.ErrorPresenter,
				data.fields.Printer,
				data.args.writer,
			)
		})
	}
}
