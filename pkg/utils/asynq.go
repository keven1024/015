package utils

import (
	"github.com/hibiken/asynq"
)

var (
	queueClient    *asynq.Client
	queueInspector *asynq.Inspector
)

func InitAsynq() error {
	opt, err := asynq.ParseRedisURI(GetEnv("redis.url"))
	if err != nil {
		return err
	}
	queueClient = asynq.NewClient(opt)
	queueInspector = asynq.NewInspector(opt)
	return nil
}

func GetQueueClient() *asynq.Client {
	return queueClient
}

func GetQueueInspector() *asynq.Inspector {
	return queueInspector
}

func RedisURI2AsynqOpt(uri string) asynq.RedisConnOpt {
	opt, err := asynq.ParseRedisURI(uri)
	if err != nil {
		panic(err)
	}
	return opt
}
