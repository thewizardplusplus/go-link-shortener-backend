package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

func TestLinkCreatingHandler_ServeHTTP(test *testing.T) {
	type fields struct {
		LinkCreator    LinkCreator
		LinkPresenter  LinkPresenter
		ErrorPresenter ErrorPresenter
	}
	type args struct {
		request *http.Request
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success",
			fields: fields{
				LinkCreator: func() LinkCreator {
					creator := new(MockLinkCreator)
					creator.
						On("CreateLink", "url").
						Return(entities.Link{Code: "code", URL: "url"}, nil)

					return creator
				}(),
				LinkPresenter: func() LinkPresenter {
					presenter := new(MockLinkPresenter)
					presenter.On(
						"PresentLink",
						mock.MatchedBy(func(writer http.ResponseWriter) bool { return true }),
						entities.Link{Code: "code", URL: "url"},
					)

					return presenter
				}(),
				ErrorPresenter: new(MockErrorPresenter),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodPost,
					"http://example.com/",
					bytes.NewBufferString(`{"URL":"url"}`),
				),
			},
		},
		{
			name: "error with decoding",
			fields: fields{
				LinkCreator:   new(MockLinkCreator),
				LinkPresenter: new(MockLinkPresenter),
				ErrorPresenter: func() ErrorPresenter {
					presenter := new(MockErrorPresenter)
					presenter.On(
						"PresentError",
						mock.MatchedBy(func(writer http.ResponseWriter) bool { return true }),
						http.StatusBadRequest,
						mock.MatchedBy(func(err error) bool { return true }),
					)

					return presenter
				}(),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodPost,
					"http://example.com/",
					bytes.NewBufferString("incorrect"),
				),
			},
		},
		{
			name: "error with creating",
			fields: fields{
				LinkCreator: func() LinkCreator {
					creator := new(MockLinkCreator)
					creator.On("CreateLink", "url").Return(entities.Link{}, iotest.ErrTimeout)

					return creator
				}(),
				LinkPresenter: new(MockLinkPresenter),
				ErrorPresenter: func() ErrorPresenter {
					presenter := new(MockErrorPresenter)
					presenter.On(
						"PresentError",
						mock.MatchedBy(func(writer http.ResponseWriter) bool { return true }),
						http.StatusInternalServerError,
						mock.MatchedBy(func(err error) bool { return true }),
					)

					return presenter
				}(),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodPost,
					"http://example.com/",
					bytes.NewBufferString(`{"URL":"url"}`),
				),
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			writer := httptest.NewRecorder()
			handler := LinkCreatingHandler{
				LinkCreator:    data.fields.LinkCreator,
				LinkPresenter:  data.fields.LinkPresenter,
				ErrorPresenter: data.fields.ErrorPresenter,
			}
			handler.ServeHTTP(writer, data.args.request)

			response := writer.Result()
			responseBody, _ := ioutil.ReadAll(response.Body)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.LinkCreator,
				data.fields.LinkPresenter,
				data.fields.ErrorPresenter,
			)
			assert.Empty(test, responseBody)
		})
	}
}
