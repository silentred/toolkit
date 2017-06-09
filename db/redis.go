package db

import (
	"log"

	"github.com/silentred/kassadin"
	redis "gopkg.in/redis.v5"
)

// NewRedisClient get a redis client
func NewRedisClient(redisInfo kassadin.RedisInstance) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisInfo.Address(),
		DB:       redisInfo.Db,
		Password: redisInfo.Pwd, // no password set
	})
	if err := client.Ping().Err(); err != nil {
		log.Fatal(err)
	}
	return client
}
