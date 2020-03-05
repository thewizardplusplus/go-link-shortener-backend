package generators

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUUIDGenerator_GenerateCode(test *testing.T) {
	for _, data := range []struct {
		name            string
		prepare         func()
		restore         func()
		wantCodePattern *regexp.Regexp
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			prepare: func() {},
			restore: func() {},
			wantCodePattern: regexp.MustCompile(
				`(?i)^[\da-f]{8}(-[\da-f]{4}){3}-[\da-f]{12}$`,
			),
			wantErr: assert.NoError,
		},
		{
			name:            "error",
			prepare:         func() { uuid.SetRand(bytes.NewReader(nil)) },
			restore:         func() { uuid.SetRand(nil) },
			wantCodePattern: regexp.MustCompile(`^$`),
			wantErr:         assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			data.prepare()
			defer data.restore()

			var generator UUIDGenerator
			gotCode, gotErr := generator.GenerateCode()

			assert.Regexp(test, data.wantCodePattern, gotCode)
			data.wantErr(test, gotErr)
		})
	}
}
