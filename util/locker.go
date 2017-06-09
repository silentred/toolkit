package util

import (
	"time"

	"gopkg.in/redis.v5"
)

type Locker interface {
	Lock(string, int) bool
	Unlock(string) bool
}

type RedisLocker struct {
	cli           *redis.Client
	defaultExpire int
}

func NewRedisLocker(cli *redis.Client, expire int) *RedisLocker {
	return &RedisLocker{cli, expire}
}

func (l *RedisLocker) Lock(key string, expireSecond int) bool {
	if l.cli != nil {
		if expireSecond == 0 {
			expireSecond = l.defaultExpire
		}
		now := time.Now().Unix()
		err := l.cli.SetNX(key, now, time.Duration(expireSecond)*time.Second).Err()

		return err == nil
	}

	return false
}

func (l *RedisLocker) Unlock(key string) bool {
	if l.cli != nil {
		err := l.cli.Del(key).Err()
		return err == nil
	}

	return false
}
