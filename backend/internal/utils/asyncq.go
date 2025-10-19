package utils

import "github.com/hibiken/asynq"

func GetQueueClient() *asynq.Client {
	opt := RedisURI2AsynqOpt(GetEnv("redis.url"))
	return asynq.NewClient(opt)
}

func GetQueueInspector() *asynq.Inspector {
	opt := RedisURI2AsynqOpt(GetEnv("redis.url"))
	return asynq.NewInspector(opt)
}

func RedisURI2AsynqOpt(uri string) asynq.RedisConnOpt {
	opt, err := asynq.ParseRedisURI(GetEnv("redis.url"))
	if err != nil {
		panic(err)
	}
	return opt
}
