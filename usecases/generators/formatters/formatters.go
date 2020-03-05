package formatters

import (
	"math/big"
	"strconv"
)

// InBase10 ...
func InBase10(code uint64) string {
	return strconv.FormatUint(code, 10)
}

// InBase62 ...
func InBase62(code uint64) string {
	var wrappedCode big.Int
	wrappedCode.SetUint64(code)

	return wrappedCode.Text(62)
}
