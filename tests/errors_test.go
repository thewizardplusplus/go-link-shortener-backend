// +build integration

package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/caarlos0/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-link-shortener/gateways/presenters"
)

func TestErrors(test *testing.T) {
	type options struct {
		ServerAddress string `env:"SERVER_ADDRESS" envDefault:"http://localhost:8080"`
	}

	var opts options
	err := env.Parse(&opts)
	require.NoError(test, err)

	for _, data := range []struct {
		name       string
		request    *http.Request
		wantStatus int
	}{
		{
			name: "unknown endpoint",
			request: func() *http.Request {
				request, _ := http.NewRequest(
					http.MethodGet,
					opts.ServerAddress+"/api/v1/incorrect",
					nil,
				)
				return request
			}(),
			wantStatus: http.StatusNotFound,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			response, err := http.DefaultClient.Do(data.request)
			require.NoError(test, err)
			defer response.Body.Close()

			var error presenters.ErrorResponse
			json.NewDecoder(response.Body).Decode(&error)

			assert.Equal(test, data.wantStatus, response.StatusCode)
			assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
			assert.NotEmpty(test, error.Error)
		})
	}
}
