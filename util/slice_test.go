package util

import "testing"
import "github.com/stretchr/testify/assert"

func TestInSlice(t *testing.T) {
	var s1 = []int{1, 2, 3, 4}
	var n = 1

	assert.True(t, InSliceInt(n, s1))
}

func TestDiff(t *testing.T) {
	var s1 = []int{1, 2, 3, 4}
	var s2 = []int{3, 4}
	s3 := SliceIntDiff(s1, s2)
	assert.True(t, len(s3) > 0)
	t.Log(s3)
}
