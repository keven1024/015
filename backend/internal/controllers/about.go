package controllers

import (
	"backend/internal/models"
	"backend/internal/utils"
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

func GetAbout(c echo.Context) error {
	maxStorageSize, err := utils.GetFileSize(utils.GetEnv("upload.maximum"))
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
		"bg_url":  utils.GetEnv("about.bg_url"),
		"content": utils.GetEnvMapString("about.content"),
		"email":   utils.GetEnv("about.email"),
		"name":    utils.GetEnv("about.name"),
		"url":     utils.GetEnv("about.url"),
		"avatar":  utils.GetEnv("about.avatar"),
		"file": map[string]any{
			"maximun": maxStorageSize,
			"current": currentFileSize,
		},
	})
}
