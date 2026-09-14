package main

import (
	"log"
	"pkg/geoip"
	"pkg/i18n"
	"pkg/utils"
	"worker/internal"

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
	utils.LoadConfig()
	// redis
	if err := utils.InitRedis(); err != nil {
		logger.Fatal("redis init failed", zap.Error(err))
		panic(err)
	}
	if err := utils.InitAsynq(); err != nil {
		logger.Fatal("asynq init failed", zap.Error(err))
		panic(err)
	}

	if err := i18n.Init(); err != nil {
		log.Fatalf("failed to init i18n: %v", err)
	}
	if err := geoip.Init(); err != nil {
		log.Fatalf("failed to init geoip: %v", err)
	}

	cfg := utils.GetConfig()
	redisConn := utils.RedisURI2AsynqOpt(cfg.Redis.URL)

	srv := asynq.NewServer(
		redisConn,
		asynq.Config{Concurrency: cast.ToInt(utils.GetEnvWithDefault("worker.concurrency", "4"))},
	)

	mux := asynq.NewServeMux()
	for _, middleware := range internal.Middlewares {
		mux.Use(middleware)
	}
	for _, route := range internal.Routes {
		mux.HandleFunc(route.Path, route.Handler)
	}

	scheduler := asynq.NewScheduler(redisConn, nil)
	if _, err := scheduler.Register("0 3 * * *", asynq.NewTask("file:janitor", nil)); err != nil {
		log.Fatalf("could not register scheduler: %v", err)
	}
	if err := scheduler.Start(); err != nil {
		log.Fatalf("could not start scheduler: %v", err)
	}
	defer scheduler.Shutdown()

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
