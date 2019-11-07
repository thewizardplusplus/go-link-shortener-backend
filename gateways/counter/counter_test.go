// +build integration

package counter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounter_NextCountChunk(test *testing.T) {
	type fields struct {
		makeClient func(test *testing.T) Client
		name       string
	}

	for _, data := range []struct {
		name    string
		prepare func(test *testing.T, counter Counter) (preparedData interface{})
		fields  fields
		wantErr assert.ErrorAssertionFunc
		check   func(test *testing.T, preparedData interface{}, gotChunk uint64)
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			client := data.fields.makeClient(test)
			counter := Counter{
				Client: client,
				Name:   data.fields.name,
			}
			preparedData := data.prepare(test, counter)
			gotChunk, gotErr := counter.NextCountChunk()

			data.wantErr(test, gotErr)
			data.check(test, preparedData, gotChunk)
		})
	}
}
