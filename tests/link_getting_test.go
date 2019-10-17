// +build integration

package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-link-shortener/entities"
	"github.com/thewizardplusplus/go-link-shortener/gateways/presenters"
)

func TestLinkGetting(test *testing.T) {
	for _, data := range []struct {
		name       string
		prepare    func(test *testing.T)
		request    *http.Request
		response   interface{}
		wantStatus int
		wantCode   string
		wantURL    string
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			data.prepare(test)

			response, err := http.DefaultClient.Do(data.request)
			require.NoError(test, err)
			defer response.Body.Close()

			json.NewDecoder(response.Body).Decode(data.response)

			assert.Equal(test, data.wantStatus, response.StatusCode)
			assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
			switch response := data.response.(type) {
			case *entities.Link:
				assert.Equal(test, data.wantCode, response.Code)
				assert.Equal(test, data.wantURL, response.URL)
			case *presenters.ErrorResponse:
				assert.NotEmpty(test, response.Error)
			default:
				test.Errorf("incorrect response type: %T", response)
			}
		})
	}
}
