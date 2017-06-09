package util

// import "testing"
// import "github.com/stretchr/testify/assert"

// func Test_MemCache(t *testing.T) {
// 	c := NewMemCache(10)
// 	tests := []struct {
// 		opt    string
// 		key    string
// 		obj    interface{}
// 		expect bool
// 	}{
// 		{"get", "aaa", "", false},
// 		{"set", "aaa", "bbb", true},
// 		{"get", "aaa", "bbb", true},
// 		{"del", "aaa", "", true},
// 		{"set", "aaa", 123, true},
// 	}

// 	for _, test := range tests {
// 		switch test.opt {
// 		case "get":
// 			obj, exists := c.Get(test.key)
// 			if exists {
// 				assert.Equal(t, test.obj, obj)
// 			}
// 			assert.Equal(t, test.expect, exists)
// 		case "set":
// 			res := c.Set(test.key, test.obj)
// 			assert.Equal(t, test.expect, res)
// 		case "del":
// 			res := c.Del(test.key)
// 			assert.Equal(t, test.expect, res)
// 		}
// 	}
// }

// func Test_TryCache(t *testing.T) {
// 	c := NewMemCache(10)
// 	obj := struct {
// 		name string
// 		age  int
// 	}{"jason", 999}

// 	tests := []struct {
// 		key    string
// 		call   Callable
// 		expect interface{}
// 	}{
// 		{"aaa", func() interface{} { return "bbb" }, "bbb"},
// 		{"bbb", func() interface{} { return 1 }, 1},
// 		{"ccc", func() interface{} { return 1.0 }, 1.0},
// 		{"ddd", func() interface{} { return true }, true},
// 		{"eee", func() interface{} { return obj }, obj},
// 		{"aaa", func() interface{} { return 123 }, "bbb"},
// 	}

// 	for _, test := range tests {
// 		res := TryCache(c, test.key, test.call)
// 		assert.Equal(t, test.expect, res)
// 	}
// }
