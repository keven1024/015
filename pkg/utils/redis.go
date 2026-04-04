package utils

import (
	"context"
	"sync"

	"github.com/redis/rueidis"
)

var (
	rdb       rueidis.Client
	ctx       = context.Background()
	onceRedis sync.Once
)

func InitRedis() rueidis.Client {
	opt, err := rueidis.ParseURL(GetEnv("redis.url"))
	if err != nil {
		panic(err)
	}
	client, err := rueidis.NewClient(opt)
	if err != nil {
		panic(err)
	}
	return client
}

func GetRedisClient() (rueidis.Client, context.Context) {
	onceRedis.Do(func() {
		InitRedis()
	})
	return rdb, ctx
}
