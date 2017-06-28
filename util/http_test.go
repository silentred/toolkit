package util

import "testing"
import "github.com/stretchr/testify/assert"

func TestHTTP(t *testing.T) {
	client := NewHTTPClient(60, nil)
	req, _ := NewHTTPReqeust("GET", "http://baidu.com", nil, nil, nil)
	b, code, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.True(t, len(b) > 0)
}

func TestParallelHTTP(t *testing.T) {
	client := NewHTTPClient(60, nil)
	req, _ := NewHTTPReqeust("GET", "http://baidu.com", nil, nil, nil)
	resps := client.DoParallel(req, req)
	for _, item := range resps {
		assert.NotNil(t, item)
		assert.NoError(t, item.Err)
		assert.True(t, item.Response.StatusCode == 200)
	}

}
