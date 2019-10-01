// +build integration

package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

func TestLinkSetter_SetLink(test *testing.T) {
	type fields struct {
		KeyExtractor KeyExtractor
		Client       Client
		Expiration   time.Duration
	}
	type args struct {
		link entities.Link
	}

	for _, data := range []struct {
		name    string
		fields  fields
		prepare func(test *testing.T, client Client)
		args    args
		wantErr assert.ErrorAssertionFunc
		check   func(test *testing.T, client Client)
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			data.prepare(test, data.fields.Client)

			cache := LinkSetter{
				KeyExtractor: data.fields.KeyExtractor,
				Client:       data.fields.Client,
				Expiration:   data.fields.Expiration,
			}
			gotErr := cache.SetLink(data.args.link)

			data.wantErr(test, gotErr)
			data.check(test, data.fields.Client)
		})
	}
}
