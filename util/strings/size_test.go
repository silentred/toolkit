package strings

import "testing"
import "github.com/stretchr/testify/assert"

func Test_ParseByteSize(t *testing.T) {
	tests := []struct {
		size   string
		result int
	}{
		{"5KB", 5 << 10},
		{"1MB", 1 << 20},
		{"1GB", 1 << 30},
		{"1K", 1 << 10},
		{"1M", 1 << 20},
		{"1G", 1 << 30},
	}

	for _, test := range tests {
		res, err := ParseByteSize(test.size)
		if assert.NoError(t, err) {
			assert.Equal(t, test.result, res)
		}
	}
}
