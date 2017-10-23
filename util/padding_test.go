package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPadding(t *testing.T) {
	tbl := []struct {
		input  string
		length int
		expect string
	}{
		{"0", 3, "000"},
		{"01", 3, "001"},
		{"9", 3, "009"},
		{"099", 3, "099"},
		{"12", 3, "012"},
	}

	for _, item := range tbl {
		ret := LeftPadding(item.input, "0", item.length)
		assert.Equal(t, item.expect, ret)
	}
}
