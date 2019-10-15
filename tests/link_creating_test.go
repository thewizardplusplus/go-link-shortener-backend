// +build integration

package tests

import (
	"encoding/json"
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

func TestLinkCreating(test *testing.T) {
	for _, data := range []struct {
		name            string
		prepare         func(test *testing.T)
		request         *http.Request
		wantStatus      int
		wantCodePattern *regexp.Regexp
		wantURL         string
		check           func(test *testing.T, expectedLink entities.Link)
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			data.prepare(test)

			response, err := http.DefaultClient.Do(data.request)
			require.NoError(test, err)
			defer response.Body.Close()

			var link entities.Link
			json.NewDecoder(response.Body).Decode(&link)

			assert.Equal(test, data.wantStatus, response.StatusCode)
			assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
			assert.Regexp(test, data.wantCodePattern, link.Code)
			assert.Equal(test, data.wantURL, link.URL)
			data.check(test, link)
		})
	}
}
