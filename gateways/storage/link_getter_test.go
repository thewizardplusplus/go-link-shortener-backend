// +build integration

package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

func TestLinkGetter_GetLink(test *testing.T) {
	type fields struct {
		makeClient func(test *testing.T) Client
		database   string
		collection string
		keyField   string
	}
	type args struct {
		query string
	}

	for _, data := range []struct {
		name     string
		fields   fields
		prepare  func(test *testing.T, getter LinkGetter)
		args     args
		wantLink entities.Link
		wantErr  assert.ErrorAssertionFunc
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			client := data.fields.makeClient(test)
			getter := LinkGetter{
				Client:     client,
				Database:   data.fields.database,
				Collection: data.fields.collection,
				KeyField:   data.fields.keyField,
			}
			data.prepare(test, getter)

			gotLink, gotErr := getter.GetLink(data.args.query)

			assert.Equal(test, data.wantLink, gotLink)
			data.wantErr(test, gotErr)
		})
	}
}
