package controllers

import (
	"backend/internal/utils"

	"github.com/labstack/echo/v4"
)

func GetAbout(c echo.Context) error {

	return utils.HTTPSuccessHandler(c, map[string]any{
		"bg_url":  utils.GetEnv("about.bg_url"),
		"content": utils.GetEnvMapString("about.content"),
		"email":   utils.GetEnv("about.email"),
		"name":    utils.GetEnv("about.name"),
		"url":     utils.GetEnv("about.url"),
		"avatar":  utils.GetEnv("about.avatar"),
	})
}
