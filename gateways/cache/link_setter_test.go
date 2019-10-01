// +build integration

package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		{
			name: "success",
			fields: fields{
				KeyExtractor: func(link entities.Link) string { return "key" },
				Client:       NewClient(address),
				Expiration:   time.Hour,
			},
			prepare: func(test *testing.T, client Client) {
				err := client.innerClient.Del("key").Err()
				require.NoError(test, err)
			},
			args: args{
				link: entities.Link{Code: "code", URL: "url"},
			},
			wantErr: assert.NoError,
			check: func(test *testing.T, client Client) {
				data, err := client.innerClient.Get("key").Result()
				require.NoError(test, err)

				duration, err := client.innerClient.TTL("key").Result()
				require.NoError(test, err)

				assert.Equal(test, `{"Code":"code","URL":"url"}`, data)
				assert.InDelta(test, time.Hour, duration, float64(10*time.Second))
			},
		},
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
