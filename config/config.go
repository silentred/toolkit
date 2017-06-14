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

type AppConfig struct {
	Name  string
	Mode  string
	Port  int
	Sess  SessionConfig
	Log   LogConfig
	Mysql MysqlConfig
	Redis RedisConfig
}

type SessionConfig struct {
	Providor  string
	StorePath string
	Enable    bool
}

type LogConfig struct {
	Name         string
	Providor     string
	LogPath      string
	RotateMode   string
	RotateLimit  string
	Suffix       string
	RotateEnable bool
}

type MysqlConfig struct {
	Instances []MysqlInstance
	InitMySQL bool
	Ping      bool
}

type MysqlInstance struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Pwd      string `json:"password"`
	Db       string `json:"db"`
	ReadOnly bool   `json:"read_only"`
}

func (inst MysqlInstance) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", inst.User, inst.Pwd, inst.Host, inst.Port, inst.Db)
}

type RedisConfig struct {
	InitRedis bool
	Ping      bool
	RedisInstance
}

type RedisInstance struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Pwd  string `json:"password"`
	Port int    `json:"port"`
	Db   int    `json:"database"`
}

func (inst RedisInstance) Address() string {
	return fmt.Sprintf("%s:%d", inst.Host, inst.Port)
}
