package counter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(test *testing.T) {
	type args struct {
		url string
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
			gotClient, gotErr := NewClient(data.args.url)

			data.wantClient(test, gotClient.innerClient)
			data.wantErr(test, gotErr)
		})
	}
}
