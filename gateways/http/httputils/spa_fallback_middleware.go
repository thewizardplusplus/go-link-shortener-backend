package httputils

import (
	"net/http"

	"github.com/golang/gddo/httputil/header"
)

func isStaticAssetRequest(request *http.Request) bool {
	if request.Method != http.MethodGet {
		return false
	}

	for _, spec := range header.ParseAccept(request.Header, "Accept") {
		if spec.Value == "text/html" {
			return true
		}
	}

	return false
}
