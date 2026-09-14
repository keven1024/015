package controllers

import (
	"backend/internal/utils"
	"pkg/models"
	u "pkg/utils"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func GetConfig(c *echo.Context) error {
	cfg := u.GetConfig()

	return utils.HTTPSuccessHandler(c, map[string]any{
		"site_title":     cfg.Site.Title,
		"site_desc":      cfg.Site.Desc,
		"site_url":       cfg.Site.URL,
		"site_icon":      cfg.Site.Icon,
		"site_bg_url":    cfg.Site.BgURL,
		"site_enable_bg": lo.FromPtr(cfg.Site.EnableBg),
		"version":        u.GetEnvWithDefault("VERSION", "dev"),
		"build_time":     cast.ToInt(u.GetEnvWithDefault("BUILD_TIME", cast.ToString(time.Now().Unix()))),
		"features": lo.FilterMap(lo.Entries(cfg.Features), func(e lo.Entry[models.FeatureConfigKey, any], _ int) (string, bool) {
			return string(e.Key), cast.ToBool(cast.ToStringMap(e.Value)["enabled"])
		}),
	})
}
