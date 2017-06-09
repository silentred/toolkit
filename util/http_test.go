package util

import "testing"
import "github.com/stretchr/testify/assert"

func TestHTTP(t *testing.T) {
	client := NewHTTPClient(60, nil)
	req, _ := NewHTTPReqeust("GET", "http://baidu.com", nil, nil, nil)
	_, code, err := client.Get(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
}
