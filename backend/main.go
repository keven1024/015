package main

import (
	"fmt"
	"pkg/utils"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

func main() {
	// 日志
	var logger *zap.Logger
	if utils.GetEnvWithDefault("node.env", "production") == "production" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync() //nolint:errcheck
	zap.ReplaceGlobals(logger)
	// redis
	if err := utils.InitRedis(); err != nil {
		logger.Fatal("redis init failed", zap.Error(err))
		panic(err)
	}
	if err := utils.InitAsynq(); err != nil {
		logger.Fatal("asynq init failed", zap.Error(err))
		panic(err)
	}

	e := echo.New()
	for _, middleware := range middlewares {
		e.Use(middleware())
	}

	for _, route := range routes {
		e.Match(route.Method, route.Path, route.Handler)
	}
	if err := e.Start(fmt.Sprintf(":%s", utils.GetEnvWithDefault("api.port", "5001"))); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}
