package controllers

import (
	"backend/internal/utils"
	"encoding/json"
	"pkg/models"
	u "pkg/utils"

	"github.com/labstack/echo/v5"
	"github.com/samber/lo"
)

func GetAbout(c *echo.Context) error {
	maxStorageSize, err := u.GetFileSize(u.GetEnv("upload.maximum"))
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	fileInfoMap, err := models.GetRedisFileInfoAll()
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	currentFileSize := lo.Reduce(lo.Values(fileInfoMap), func(agg int64, item string, _ int) int64 {
		var fileInfo models.RedisFileInfo
		err := json.Unmarshal([]byte(item), &fileInfo)
		if err != nil {
			return agg
		}
		return agg + fileInfo.FileSize
	}, 0)

	return utils.HTTPSuccessHandler(c, map[string]any{
		"bg_url":  u.GetEnv("about.bg_url"),
		"content": u.GetEnvMap("about.content"),
		"email":   u.GetEnv("about.email"),
		"name":    u.GetEnv("about.name"),
		"url":     u.GetEnv("about.url"),
		"avatar":  u.GetEnv("about.avatar"),
		"file": map[string]any{
			"maximun": maxStorageSize,
			"current": currentFileSize,
		},
	})
}
