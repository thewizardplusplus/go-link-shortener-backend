// +build integration

package tests

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/caarlos0/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	staticFileContent = `
		<!DOCTYPE html>
		<html lang="en">
		<title>Minimal HTML5 Document</title>
	`
)

func TestStaticFile(test *testing.T) {
	// nolint: lll
	type options struct {
		ServerAddress    string `env:"SERVER_ADDRESS" envDefault:"http://localhost:8080"`
		ServerStaticPath string `env:"SERVER_STATIC_PATH" envDefault:"../static"`
	}

	var opts options
	err := env.Parse(&opts)
	require.NoError(test, err)

	for _, data := range []struct {
		name            string
		prepare         func(test *testing.T)
		restore         func(test *testing.T)
		request         *http.Request
		wantStatus      int
		wantContentType string
		wantBody        string
	}{
		{
			name: "with a file/without the SPA fallback",
			prepare: func(test *testing.T) {
				path := filepath.Join(opts.ServerStaticPath, "page.html")
				err := ioutil.WriteFile(path, []byte(staticFileContent), 0644)
				require.NoError(test, err)
			},
			restore: func(test *testing.T) {
				path := filepath.Join(opts.ServerStaticPath, "page.html")
				err := os.Remove(path)
				require.NoError(test, err)
			},
			request: func() *http.Request {
				request, _ := http.NewRequest(
					http.MethodGet,
					opts.ServerAddress+"/page.html",
					nil,
				)
				return request
			}(),
			wantStatus:      http.StatusOK,
			wantContentType: "text/html; charset=utf-8",
			wantBody:        staticFileContent,
		},
		{
			name: "with a file/with the SPA fallback",
			prepare: func(test *testing.T) {
				path := filepath.Join(opts.ServerStaticPath, "index.html")
				err := ioutil.WriteFile(path, []byte(staticFileContent), 0644)
				require.NoError(test, err)
			},
			restore: func(test *testing.T) {
				path := filepath.Join(opts.ServerStaticPath, "index.html")
				err := os.Remove(path)
				require.NoError(test, err)
			},
			request: func() *http.Request {
				request, _ := http.NewRequest(
					http.MethodGet,
					opts.ServerAddress+"/page.html",
					nil,
				)
				request.Header.Set("Accept", "text/html")

				return request
			}(),
			wantStatus:      http.StatusOK,
			wantContentType: "text/html; charset=utf-8",
			wantBody:        staticFileContent,
		},
		{
			name: "without a file",
			prepare: func(test *testing.T) {
				path := filepath.Join(opts.ServerStaticPath, "page.html")
				if err := os.Remove(path); !os.IsNotExist(err) {
					require.NoError(test, err)
				}
			},
			restore: func(test *testing.T) {},
			request: func() *http.Request {
				request, _ := http.NewRequest(
					http.MethodGet,
					opts.ServerAddress+"/page.html",
					nil,
				)
				return request
			}(),
			wantStatus:      http.StatusNotFound,
			wantContentType: "text/plain; charset=utf-8",
			wantBody:        "404 page not found\n",
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			data.prepare(test)
			defer data.restore(test)

			response, err := http.DefaultClient.Do(data.request)
			require.NoError(test, err)
			defer response.Body.Close()

			responseBody, err := ioutil.ReadAll(response.Body)
			require.NoError(test, err)

			assert.Equal(test, data.wantStatus, response.StatusCode)
			assert.Equal(test, data.wantContentType, response.Header.Get("Content-Type"))
			assert.Equal(test, data.wantBody, string(responseBody))
		})
	}
}
