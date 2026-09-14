package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"pkg/utils"
	"syscall"
	"time"

	"backend/internal"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

func main() {
	os.Exit(run())
}

func run() int {
	// 日志
	var logger *zap.Logger
	if utils.GetEnvWithDefault("node.env", "production") == "production" {
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

	e := echo.New()
	// Without this, RealIP() returns 127.0.0.1 (Caddy's loopback proxy), not the real client.
	e.IPExtractor = echo.ExtractIPFromXFFHeader(echo.TrustLoopback(true))
	e.Filesystem = os.DirFS(utils.GetConfig().Upload.Path)
	for _, middleware := range internal.Middlewares {
		e.Use(middleware())
	}

	for _, route := range internal.Routes {
		e.Match(route.Method, route.Path, route.Handler)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := (echo.StartConfig{
		Address:         fmt.Sprintf(":%s", utils.GetEnvWithDefault("api.port", "5001")),
		GracefulTimeout: 30 * time.Second,
	}).Start(ctx, e); err != nil {
		logger.Error("server failed", zap.Error(err))
		return 1
	}
	return 0
}
