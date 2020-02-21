// +build integration

package tests

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/caarlos0/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

			responseBody, err := ioutil.ReadAll(response.Body)
			require.NoError(test, err)

			assert.Equal(test, data.wantStatus, response.StatusCode)
			assert.Equal(
				test,
				"text/plain; charset=utf-8",
				response.Header.Get("Content-Type"),
			)
			assert.NotEmpty(test, responseBody)
		})
	}
}
