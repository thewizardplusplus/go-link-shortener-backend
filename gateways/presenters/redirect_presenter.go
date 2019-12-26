package presenters

import (
	"html/template"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

// RedirectPresenter ...
type RedirectPresenter struct {
	ErrorURL string
	Printer  Printer
}

// PresentLink ...
func (presenter RedirectPresenter) PresentLink(
	writer http.ResponseWriter,
	link entities.Link,
) error {
	if err := redirect(writer, link.URL, http.StatusMovedPermanently); err != nil {
		return errors.Wrap(err, "unable to redirect to the link")
	}

	return nil
}

// PresentError ...
func (presenter RedirectPresenter) PresentError(
	writer http.ResponseWriter,
	statusCode int,
	err error,
) error {
	err2 := redirect(writer, presenter.ErrorURL, http.StatusFound)
	if err2 != nil {
		return errors.Wrap(err2, "unable to redirect to the error")
	}

	presenter.Printer.Printf("redirect because of the error: %v", err)
	return nil
}

var (
	// it's based on:
	// * https://www.sitepoint.com/a-minimal-html-document-html5-edition/
	// * https://developer.mozilla.org/en-US/docs/Mozilla/Mobile/Viewport_meta_tag
	responseTemplate = template.Must(template.New("redirect").Parse(`
		<!DOCTYPE html>
		<html lang="en">
			<head>
				<meta charset="utf-8" />
				<meta name="viewport" content="width=device-width, initial-scale=1" />

				<title>Redirect</title>
			</head>
			<body>
				<p>{{ .StatusText }}: <a href="{{ .URL }}">{{ .URL }}</a></p>
			</body>
		</html>
	`))
)

// we use our implementation of redirection
// because the http.Redirect function doesn't return any errors;
//
// errors with writing to the http.ResponseWriter is important to handle,
// see for details: https://stackoverflow.com/a/43976633
func redirect(writer http.ResponseWriter, url string, statusCode int) error {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.Header().Set("Location", url)
	writer.WriteHeader(statusCode)

	// we respond with a short HTML body by recommendations from RFC 7231
	err := responseTemplate.Execute(writer, struct {
		StatusText string
		URL        string
	}{
		StatusText: http.StatusText(statusCode),
		URL:        url,
	})
	if err != nil {
		return errors.Wrap(err, "unable to write the data")
	}

	return nil
}
