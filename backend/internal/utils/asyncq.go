package utils

import "github.com/hibiken/asynq"

func GetQueueClient() *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: GetEnv("REDIS_URL")})
}

func GetQueueInspector() *asynq.Inspector {
	return asynq.NewInspector(asynq.RedisClientOpt{Addr: GetEnv("REDIS_URL")})
}
