package db

import (
	"log"

	"github.com/silentred/toolkit/config"
	redis "gopkg.in/redis.v5"
)

// NewRedisClient get a redis client
func NewRedisClient(redisInfo config.RedisInstance) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisInfo.Address(),
		DB:       redisInfo.Db,
		Password: redisInfo.Pwd, // no password set
	})

	if redisInfo.Ping {
		if err := client.Ping().Err(); err != nil {
			log.Fatal(err)
		}
	}

	return client
}
