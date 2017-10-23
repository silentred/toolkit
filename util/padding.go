package util

import (
	"strings"
)

func LeftPadding(str, pad string, length int) string {
	if len(str) < length {
		diff := length - len(str)
		padLen := len(pad)
		times := diff/padLen + 1
		left := strings.Repeat(pad, times)
		tmp := left + str
		idx := len(tmp) - length
		return tmp[idx:]
	}

	return str
}
