package utils

import "github.com/hibiken/asynq"

func RedisURI2AsynqOpt(uri string) asynq.RedisConnOpt {
	opt, err := asynq.ParseRedisURI(GetEnv("redis.url"))
	if err != nil {
		panic(err)
	}
	return opt
}
