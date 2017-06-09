package db

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func TestConfig(t *testing.T) {
	config := []byte(`
[[products]]
name = "Hammer"
sku = 738594937

[[products]]
name = "Nail"
sku = 284758393`)

	viper.SetConfigType("toml")
	viper.ReadConfig(bytes.NewBuffer(config))

	obj := viper.Get("products")
	fmt.Printf("%#v", obj)

}
