package main

import (
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
		logger.Fatal("failed to init i18n:", zap.Error(err))
	}
	if err := geoip.Init(); err != nil {
		logger.Fatal("failed to init geoip", zap.Error(err))
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
	for _, schedule := range internal.Schedules {
		if _, err := scheduler.Register(schedule.Cron, schedule.Handle); err != nil {
			logger.Fatal("could not register scheduler", zap.Error(err))
		}
	}
	if err := scheduler.Start(); err != nil {
		logger.Fatal("could not start scheduler", zap.Error(err))
	}
	defer scheduler.Shutdown()

	if err := srv.Run(mux); err != nil {
		logger.Fatal("could not run server", zap.Error(err))
	}
}
