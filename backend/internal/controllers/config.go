package controllers

import (
	"backend/internal/utils"

	"github.com/labstack/echo/v4"
)

func GetConfig(c echo.Context) error {
	return utils.HTTPSuccessHandler(c, map[string]any{
		"site_title": utils.GetEnv("site_title"),
		"site_desc":  utils.GetEnv("site_desc"),
		"site_url":   utils.GetEnv("site_url"),
	})
}
