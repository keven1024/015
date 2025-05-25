package main

import (
	"log"
	"worker/internal/tasks"
	"worker/internal/utils"

	"github.com/hibiken/asynq"
	"github.com/spf13/cast"
)

func main() {
	srv := asynq.NewServer(
		utils.RedisURI2AsynqOpt(utils.GetEnv("REDIS_URL")),
		asynq.Config{Concurrency: cast.ToInt(utils.GetEnvWithDefault("WORKER_CONCURRENCY", "4"))},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("image:compress", tasks.CompressImage)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
