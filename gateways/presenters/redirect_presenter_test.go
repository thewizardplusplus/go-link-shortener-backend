package presenters

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

// indents inside constants are significant
const (
	responseAtLinkPresenting = `
		<!DOCTYPE html>
		<html lang="en">
			<head>
				<meta charset="utf-8" />
				<meta name="viewport" content="width=device-width, initial-scale=1" />

				<title>Redirect</title>
			</head>
			<body>
				<p>Moved Permanently: <a href="http://example.com/">http://example.com/</a></p>
			</body>
		</html>
	`
	responseAtErrorPresenting = `
		<!DOCTYPE html>
		<html lang="en">
			<head>
				<meta charset="utf-8" />
				<meta name="viewport" content="width=device-width, initial-scale=1" />

				<title>Redirect</title>
			</head>
			<body>
				<p>Found: <a href="/error">/error</a></p>
			</body>
		</html>
	`
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
		{
			name: "success",
			fields: fields{
				ErrorURL: "/error",
				Printer:  new(MockPrinter),
			},
			args: args{
				writer: httptest.NewRecorder(),
				link:   entities.Link{Code: "code", URL: "http://example.com/"},
			},
			wantErr: assert.NoError,
			check: func(test *testing.T, writer http.ResponseWriter) {
				response := writer.(*httptest.ResponseRecorder).Result()
				responseBody, _ := ioutil.ReadAll(response.Body)

				assert.Equal(test, http.StatusMovedPermanently, response.StatusCode)
				assert.Equal(
					test,
					"text/html; charset=utf-8",
					response.Header.Get("Content-Type"),
				)
				assert.Equal(test, "http://example.com/", response.Header.Get("Location"))
				assert.Equal(test, responseAtLinkPresenting, string(responseBody))
			},
		},
		{
			name: "error",
			fields: fields{
				ErrorURL: "/error",
				Printer:  new(MockPrinter),
			},
			args: args{
				writer: NewTimeoutResponseRecorder(),
				link:   entities.Link{Code: "code", URL: "http://example.com/"},
			},
			wantErr: assert.Error,
			check: func(test *testing.T, writer http.ResponseWriter) {
				response := writer.(TimeoutResponseRecorder).Result()
				responseBody, _ := ioutil.ReadAll(response.Body)

				assert.Equal(test, http.StatusMovedPermanently, response.StatusCode)
				assert.Equal(
					test,
					"text/html; charset=utf-8",
					response.Header.Get("Content-Type"),
				)
				assert.Equal(test, "http://example.com/", response.Header.Get("Location"))
				assert.Empty(test, responseBody)
			},
		},
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

func TestRedirectPresenter_PresentError(test *testing.T) {
	type fields struct {
		ErrorURL string
		Printer  Printer
	}
	type args struct {
		writer     http.ResponseWriter
		statusCode int
		err        error
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
		check   func(test *testing.T, writer http.ResponseWriter)
	}{
		{
			name: "success",
			fields: fields{
				ErrorURL: "/error",
				Printer: func() Printer {
					printer := new(MockPrinter)
					printer.
						On(
							"Printf",
							mock.MatchedBy(func(string) bool { return true }),
							iotest.ErrTimeout,
						).
						Return()

					return printer
				}(),
			},
			args: args{
				writer:     httptest.NewRecorder(),
				statusCode: http.StatusInternalServerError,
				err:        iotest.ErrTimeout,
			},
			wantErr: assert.NoError,
			check: func(test *testing.T, writer http.ResponseWriter) {
				response := writer.(*httptest.ResponseRecorder).Result()
				responseBody, _ := ioutil.ReadAll(response.Body)

				assert.Equal(test, http.StatusFound, response.StatusCode)
				assert.Equal(
					test,
					"text/html; charset=utf-8",
					response.Header.Get("Content-Type"),
				)
				assert.Equal(test, "/error", response.Header.Get("Location"))
				assert.Equal(test, responseAtErrorPresenting, string(responseBody))
			},
		},
		{
			name: "error",
			fields: fields{
				ErrorURL: "/error",
				Printer:  new(MockPrinter),
			},
			args: args{
				writer:     NewTimeoutResponseRecorder(),
				statusCode: http.StatusInternalServerError,
				err:        iotest.ErrTimeout,
			},
			wantErr: assert.Error,
			check: func(test *testing.T, writer http.ResponseWriter) {
				response := writer.(TimeoutResponseRecorder).Result()
				responseBody, _ := ioutil.ReadAll(response.Body)

				assert.Equal(test, http.StatusFound, response.StatusCode)
				assert.Equal(
					test,
					"text/html; charset=utf-8",
					response.Header.Get("Content-Type"),
				)
				assert.Equal(test, "/error", response.Header.Get("Location"))
				assert.Empty(test, responseBody)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			presenter := RedirectPresenter{
				ErrorURL: data.fields.ErrorURL,
				Printer:  data.fields.Printer,
			}
			gotErr := presenter.PresentError(
				data.args.writer,
				data.args.statusCode,
				data.args.err,
			)

			mock.AssertExpectationsForObjects(test, data.fields.Printer)
			data.wantErr(test, gotErr)
			data.check(test, data.args.writer)
		})
	}
}
