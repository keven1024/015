package main

import (
	"backend/internal/controllers"

	"github.com/labstack/echo/v5"
)

type Route struct {
	Method  []string
	Path    string
	Handler func(c *echo.Context) error
}

var routes = []Route{
	{Method: []string{"POST"}, Path: "/file/create", Handler: controllers.CreateUploadTask},
	{Method: []string{"POST"}, Path: "/file/slice", Handler: controllers.UploadFileSlice},
	{Method: []string{"POST"}, Path: "/file/finish", Handler: controllers.FinishUploadTask},

	{Method: []string{"GET"}, Path: "/share/:id", Handler: controllers.GetShareInfo},
	{Method: []string{"POST"}, Path: "/share", Handler: controllers.CreateShareInfo},
	{Method: []string{"DELETE"}, Path: "/share/:id", Handler: controllers.DeleteShareInfo},

	{Method: []string{"GET"}, Path: "/download", Handler: controllers.DownloadShare},
	{Method: []string{"POST"}, Path: "/download", Handler: controllers.VaildateShare},
	{Method: []string{"GET"}, Path: "/share/pickup/:code", Handler: controllers.GetShareByPickupCode},

	{Method: []string{"GET"}, Path: "/stat", Handler: controllers.GetStat},
	{Method: []string{"GET"}, Path: "/config", Handler: controllers.GetConfig},
	{Method: []string{"GET"}, Path: "/about", Handler: controllers.GetAbout},

	{Method: []string{"POST"}, Path: "/task/:type", Handler: controllers.CreateTask},
	{Method: []string{"GET"}, Path: "/task/:id", Handler: controllers.GetTask},
}
