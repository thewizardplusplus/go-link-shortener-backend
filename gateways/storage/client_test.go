package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(test *testing.T) {
	type args struct {
		uri string
	}

	for _, data := range []struct {
		name       string
		args       args
		wantClient assert.ValueAssertionFunc
		wantErr    assert.ErrorAssertionFunc
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			gotClient, gotErr := NewClient(data.args.uri)

			data.wantClient(test, gotClient.innerClient)
			data.wantErr(test, gotErr)
		})
	}
}
