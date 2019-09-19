package code

import (
	"regexp"
	"testing"

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
		// TODO: add test cases
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
