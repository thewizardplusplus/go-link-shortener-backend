// +build integration

package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

func TestLinkGetter_GetLink(test *testing.T) {
	type fields struct {
		Client Client
	}
	type args struct {
		query string
	}

	for _, data := range []struct {
		name     string
		fields   fields
		prepare  func(test *testing.T, client Client)
		args     args
		wantLink entities.Link
		wantErr  assert.ErrorAssertionFunc
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			data.prepare(test, data.fields.Client)

			cache := LinkGetter{
				Client: data.fields.Client,
			}
			gotLink, gotErr := cache.GetLink(data.args.query)

			assert.Equal(test, data.wantLink, gotLink)
			data.wantErr(test, gotErr)
		})
	}
}
