package utils

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client = InitRedis()
var ctx = context.Background()

func InitRedis() *redis.Client {
	opt, err := redis.ParseURL("redis://192.168.100.5:6379/3")
	if err != nil {
		panic(err)
	}
	return redis.NewClient(opt)
}

func GetRedisClient() (*redis.Client, context.Context) {
	return rdb, ctx
}
