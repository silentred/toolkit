package strings

import (
	"strconv"
	"strings"
	"unicode"
)

const (
	Byte  = 0
	KByte = 10
	MByte = 20
	GByte = 30
)

func ParseByteSize(size string) (int, error) {
	var moveOffset uint
	var lastDigitOffset int

	rules := []struct {
		suffix []string
		offset uint
	}{
		{[]string{"KB", "K"}, KByte},
		{[]string{"MB", "M"}, MByte},
		{[]string{"GB", "G"}, GByte},
	}

	for _, val := range rules {
		if strings.HasSuffix(size, val.suffix[0]) || strings.HasSuffix(size, val.suffix[1]) {
			moveOffset = val.offset
			break
		}
	}

	for _, s := range size {
		if unicode.IsDigit(s) {
			lastDigitOffset++
		} else {
			break
		}
	}

	val, err := strconv.Atoi(size[:lastDigitOffset])
	if err != nil {
		return 0, err
	}
	val = val << moveOffset

	return val, nil
}
