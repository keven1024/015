package main

import (
	"log"
	"worker/internal/tasks"
	"worker/internal/utils"
	"worker/middleware"

	"github.com/hibiken/asynq"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

func main() {
	// 日志
	var logger *zap.Logger
	if utils.GetEnvWithDefault("NODE_ENV", "production") == "production" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	srv := asynq.NewServer(
		utils.RedisURI2AsynqOpt(utils.GetEnv("redis.url")),
		asynq.Config{Concurrency: cast.ToInt(utils.GetEnvWithDefault("worker.concurrency", "4"))},
	)

	mux := asynq.NewServeMux()
	mux.Use(middleware.LoggerMiddleware)
	mux.HandleFunc("share:remove", tasks.RemoveShare)
	mux.HandleFunc("file:remove", tasks.RemoveFile)
	mux.HandleFunc("image:compress", tasks.CompressImage)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
