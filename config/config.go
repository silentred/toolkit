package config

import (
	"fmt"
)

const (
	ModeDev  = "dev"
	ModeProd = "prod"

	ProvidorFile   = "file"
	ProvidorStdOut = "stdout"

	RotateByDay  = "day"
	RotateBySize = "size"
)

// AppConfig for application
type AppConfig struct {
	Name  string
	Mode  string
	Port  int
	Log   LogConfig
	Mysql MysqlConfig
	Redis RedisConfig
}

type sessionConfig struct {
	Providor  string
	StorePath string
	Enable    bool
}

// LogConfig for logger
type LogConfig struct {
	Name         string
	Providor     string
	LogPath      string
	RotateMode   string
	RotateLimit  string
	Suffix       string
	RotateEnable bool
}

// MysqlConfig for MySQL
type MysqlConfig struct {
	Instances []MysqlInstance
	InitMySQL bool
	Ping      bool
}

// MysqlInstance represents a single instance of mysql server
type MysqlInstance struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Pwd      string `json:"password"`
	Db       string `json:"db"`
	Version  string `json:"version"`
	Port     int    `json:"port"`
	ReadOnly bool   `json:"read_only"`
}

func (inst MysqlInstance) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", inst.User, inst.Pwd, inst.Host, inst.Port, inst.Db)
}

// RedisConfig for redis
type RedisConfig struct {
	InitRedis bool
	Ping      bool
	RedisInstance
}

// RedisInstance represents a single instance of redis server
type RedisInstance struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Pwd  string `json:"password"`
	Port int    `json:"port"`
	Db   int    `json:"database"`
}

// Address returns the address of redis server
func (inst RedisInstance) Address() string {
	return fmt.Sprintf("%s:%d", inst.Host, inst.Port)
}
