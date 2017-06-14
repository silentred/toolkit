package service

import (
	"log"

	"github.com/silentred/toolkit/config"
	redis "gopkg.in/redis.v5"
)

// NewRedisClient get a redis client
func NewRedisClient(cfg config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisInstance.Address(),
		DB:       cfg.RedisInstance.Db,
		Password: cfg.RedisInstance.Pwd, // no password set
	})

	if cfg.Ping {
		if err := client.Ping().Err(); err != nil {
			log.Fatal(err)
		}
	}

	return client
}

func initRedis(app Application) error {
	if app.GetConfig().Redis.InitRedis {
		redis := NewRedisClient(app.GetConfig().Redis)
		app.Set("redis", redis, nil)
	}
	return nil
}
