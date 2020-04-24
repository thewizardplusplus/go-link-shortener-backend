package formatters

import (
	"math/big"
)

// InBase62 ...
func InBase62(code uint64) string {
	var wrappedCode big.Int
	wrappedCode.SetUint64(code)

	return wrappedCode.Text(62)
}
