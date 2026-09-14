package main

import (
	"context"
	"os"
	"os/signal"
	"pkg/geoip"
	"pkg/i18n"
	"pkg/utils"
	"syscall"
	"time"
	"worker/internal"

	"github.com/hibiken/asynq"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

func main() {
	os.Exit(run())
}

func run() int {
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
		logger.Error("redis init failed", zap.Error(err))
		return 1
	}
	defer utils.GetRedisClient().Close()
	if err := utils.InitAsynq(); err != nil {
		logger.Error("asynq init failed", zap.Error(err))
		return 1
	}
	defer utils.GetQueueClient().Close()    //nolint:errcheck
	defer utils.GetQueueInspector().Close() //nolint:errcheck

	if err := i18n.Init(); err != nil {
		logger.Error("failed to init i18n", zap.Error(err))
		return 1
	}
	if err := geoip.Init(); err != nil {
		logger.Error("failed to init geoip", zap.Error(err))
		return 1
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	redisConn := utils.RedisURI2AsynqOpt(utils.GetConfig().Redis.URL)

	srv := asynq.NewServer(
		redisConn,
		asynq.Config{
			Concurrency:     cast.ToInt(utils.GetEnvWithDefault("worker.concurrency", "4")),
			ShutdownTimeout: 30 * time.Second,
		},
	)

	mux := asynq.NewServeMux()
	for _, middleware := range internal.Middlewares {
		mux.Use(middleware)
	}
	for _, route := range internal.Routes {
		mux.HandleFunc(route.Path, route.Handler)
	}
	if err := srv.Start(mux); err != nil {
		logger.Error("could not start server", zap.Error(err))
		return 1
	}
	defer srv.Shutdown()

	scheduler := asynq.NewScheduler(redisConn, nil)
	for _, schedule := range internal.Schedules {
		if _, err := scheduler.Register(schedule.Cron, schedule.Handle); err != nil {
			logger.Error("could not register scheduler", zap.Error(err))
			return 1
		}
	}
	if err := scheduler.Start(); err != nil {
		logger.Error("could not start scheduler", zap.Error(err))
		return 1
	}
	defer scheduler.Shutdown()

	<-ctx.Done()
	stop()
	logger.Info("shutting down worker")
	srv.Stop()
	return 0
}
