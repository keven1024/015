package utils

import "github.com/hibiken/asynq"

func GetQueueClient() *asynq.Client {
	opt := RedisURI2AsynqOpt(GetEnv("REDIS_URL"))
	return asynq.NewClient(opt)
}

func GetQueueInspector() *asynq.Inspector {
	opt := RedisURI2AsynqOpt(GetEnv("REDIS_URL"))
	return asynq.NewInspector(opt)
}

func RedisURI2AsynqOpt(uri string) asynq.RedisConnOpt {
	opt, err := asynq.ParseRedisURI(GetEnv("REDIS_URL"))
	if err != nil {
		panic(err)
	}
	return opt
}
