package controllers

import (
	"backend/internal/utils"

	"github.com/labstack/echo/v4"
)

func GetConfig(c echo.Context) error {
	return utils.HTTPSuccessHandler(c, map[string]any{
		"site_title":  utils.GetEnvMapString("site.title"),
		"site_desc":   utils.GetEnvMapString("site.desc"),
		"site_url":    utils.GetEnv("site.url"),
		"site_icon":   utils.GetEnvWithDefault("site.icon", "/logo.png"),
		"site_bg_url": utils.GetEnvWithDefault("site.bg_url", "https://img.fudaoyuan.icu/api/1/random/?scale_min=1.5&webp=true&md=false&format=302"),
	})
}
