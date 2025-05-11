package main

import (
	"backend/internal/controllers"
	"backend/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Use(middleware.ContextMiddleware())
	e.Use(middleware.SessionMiddleware())
	e.Use(middleware.AuthMiddleware())

	// e.GET("/file/:id", controllers.GetFile)
	e.POST("/file/create", controllers.CreateUploadTask)
	e.POST("/file/slice", controllers.UploadFileSlice)
	e.POST("/file/finish", controllers.FinishUploadTask)
	e.GET("/share/:id", controllers.GetShareInfo)
	e.POST("/share", controllers.CreateShareInfo)
	e.GET("/config", controllers.GetConfig)

	e.GET("/download/:id", controllers.DownloadShare)

	e.Logger.Fatal(e.Start(":1323"))
}
