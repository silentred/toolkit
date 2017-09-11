package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := NewError(404, "not found")
	assert.Equal(t, 404, err.Code)
}
