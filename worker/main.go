package main

import (
	"log"
	"pkg/i18n"
	"pkg/utils"
	"worker/internal/tasks"
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
	defer logger.Sync() //nolint:errcheck
	zap.ReplaceGlobals(logger)

	if err := i18n.Init(); err != nil {
		log.Fatalf("failed to init i18n: %v", err)
	}

	srv := asynq.NewServer(
		utils.RedisURI2AsynqOpt(utils.GetEnv("redis.url")),
		asynq.Config{Concurrency: cast.ToInt(utils.GetEnvWithDefault("worker.concurrency", "4"))},
	)

	mux := asynq.NewServeMux()
	mux.Use(middleware.LoggerMiddleware)
	mux.HandleFunc("share:remove", tasks.RemoveShare)
	mux.HandleFunc("share:notify", tasks.ShareNotify)
	mux.HandleFunc("file:remove", tasks.RemoveFile)
	mux.HandleFunc("image:compress", tasks.CompressImage)
	mux.HandleFunc("image:convert", tasks.ConvertImage)
	mux.HandleFunc("text:translate", tasks.TranslateText)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
