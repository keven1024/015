package controllers

import (
	"backend/internal/utils"
	u "pkg/utils"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func GetConfig(c *echo.Context) error {
	featureConfig := u.GetEnvMap("features")
	features := lo.FilterMap(lo.Entries(featureConfig), func(e lo.Entry[string, any], _ int) (string, bool) {
		node, ok := e.Value.(map[string]any)
		return e.Key, ok && cast.ToBool(node["enabled"])
	})

	return utils.HTTPSuccessHandler(c, map[string]any{
		"site_title":  u.GetEnvMap("site.title"),
		"site_desc":   u.GetEnvMap("site.desc"),
		"site_url":    u.GetEnv("site.url"),
		"site_icon":   u.GetEnvWithDefault("site.icon", "/logo.png"),
		"site_bg_url": u.GetEnvWithDefault("site.bg_url", "https://img.fudaoyuan.icu/api/1/random/?scale_min=1.5&webp=true&md=false&format=302"),
		"version":     u.GetEnvWithDefault("VERSION", "dev"),
		"build_time":  cast.ToInt(u.GetEnvWithDefault("BUILD_TIME", cast.ToString(time.Now().Unix()))),
		"features":    features,
	})
}
