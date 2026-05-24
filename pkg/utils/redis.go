package utils

import (
	"github.com/redis/rueidis"
)

var rdb rueidis.Client

func InitRedis() error {
	opt, err := rueidis.ParseURL(GetEnv("redis.url"))
	if err != nil {
		return err
	}
	client, err := rueidis.NewClient(opt)
	if err != nil {
		return err
	}
	rdb = client
	return nil
}

func GetRedisClient() rueidis.Client {
	return rdb
}
