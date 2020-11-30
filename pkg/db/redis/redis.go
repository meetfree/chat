// https://godoc.org/github.com/go-redis/redis
package redis

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"im-service/conf"
)

var (
	Conn *redis.Client
	err  error
)

func Open() {
	Conn = redis.NewClient(&redis.Options{
		Addr:     conf.TomlConfig.RedisConfig(gin.Mode()).Addr,
		Password: conf.TomlConfig.RedisConfig(gin.Mode()).Password,
		DB:       0,
	})
}

func Get(key string) string {
	return Conn.Get(context.TODO(), key).Val()
}

func Set(key string, val string, exp time.Duration) error {
	return Conn.Set(context.TODO(), key, val, exp).Err()
}

func HGet(key, field string) string {
	return Conn.HGet(context.TODO(), key, field).Val()
}

func Pool() chan *redis.Client {
	pool := make(chan *redis.Client, 10)
	for i := 0; i <= 10; i++ {
		pool <- new(redis.Client)
	}
	return pool
}
