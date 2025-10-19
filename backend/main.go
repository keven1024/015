package main

import (
	"backend/internal/controllers"
	"backend/internal/utils"
	"backend/middleware"
	"fmt"

	"github.com/labstack/echo/v4"
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
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	e := echo.New()
	e.Use(middleware.ContextMiddleware())
	e.Use(middleware.SessionMiddleware())
	e.Use(middleware.AuthMiddleware())
	e.Use(middleware.RateLimiterMiddleware())
	e.Use(middleware.LoggerMiddleware())

	e.POST("/file/create", controllers.CreateUploadTask)
	e.POST("/file/slice", controllers.UploadFileSlice)
	e.POST("/file/finish", controllers.FinishUploadTask)
	e.GET("/share/:id", controllers.GetShareInfo)
	e.POST("/share", controllers.CreateShareInfo)
	e.GET("/download", controllers.DownloadShare)
	e.POST("/download", controllers.VaildateShare)
	e.GET("/share/pickup/:code", controllers.GetShareByPickupCode)

	e.POST("/image/compress", controllers.GenCompressImage)
	e.GET("/image/compress/:id", controllers.GetCompressImage)

	e.GET("/stat", controllers.GetStat)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", utils.GetEnvWithDefault("api.port", "5001"))))
}
