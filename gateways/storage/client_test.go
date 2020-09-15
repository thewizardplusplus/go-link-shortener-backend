// +build ignore

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
		{
			name:       "success",
			args:       args{"mongodb://localhost:27017"},
			wantClient: assert.NotNil,
			wantErr:    assert.NoError,
		},
		{
			name:       "error",
			args:       args{":"},
			wantClient: assert.Nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotClient, gotErr := NewClient(data.args.uri)

			data.wantClient(test, gotClient.innerClient)
			data.wantErr(test, gotErr)
		})
	}
}
