package handlers

import (
	"net/http"

	"github.com/pkg/errors"
)

// NotFoundHandler ...
type NotFoundHandler struct {
	TextErrorPresenter ErrorPresenter
	JSONErrorPresenter ErrorPresenter
}

// ServeHTTP ...
func (handler NotFoundHandler) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var errorPresenter ErrorPresenter
	if isTextAccepted(request) {
		errorPresenter = handler.TextErrorPresenter
	} else {
		errorPresenter = handler.JSONErrorPresenter
	}

	const statusCode = http.StatusNotFound
	err := errors.New("unable to find the endpoint")
	errorPresenter.PresentError(writer, request, statusCode, err)
}

// based on:
// * https://create-react-app.dev/docs/proxying-api-requests-in-development
// * https://github.com/facebook/create-react-app/blob/v3.4.0/packages/react-dev-utils/WebpackDevServerUtils.js#L423
func isTextAccepted(request *http.Request) bool {
	if request.Method != http.MethodGet {
		return false
	}

	var isTextAccepted bool
	for _, acceptingType := range request.Header["Accept"] {
		if acceptingType == "text/html" {
			isTextAccepted = true
			break
		}
	}

	return isTextAccepted
}
