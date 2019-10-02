// +build integration

package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

func TestLinkSetter_SetLink(test *testing.T) {
	type fields struct {
		makeClient func(test *testing.T) Client
		database   string
		collection string
	}
	type args struct {
		link entities.Link
	}

	for _, data := range []struct {
		name    string
		fields  fields
		prepare func(test *testing.T, setter LinkSetter)
		args    args
		wantErr assert.ErrorAssertionFunc
		check   func(test *testing.T, setter LinkSetter)
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			client := data.fields.makeClient(test)
			setter := LinkSetter{
				Client:     client,
				Database:   data.fields.database,
				Collection: data.fields.collection,
			}
			data.prepare(test, setter)

			gotErr := setter.SetLink(data.args.link)

			data.wantErr(test, gotErr)
			data.check(test, setter)
		})
	}
}
