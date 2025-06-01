package controllers

import (
	"backend/internal/models"
	"backend/internal/utils"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
)

func GetStat(c echo.Context) error {
	keys, err := models.GetRedisFileKeysAll()
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	var filesSize int64
	for _, key := range keys {
		list := strings.Split(key, "_")
		if len(list) > 1 {
			filesSize += cast.ToInt64(list[1])
		}
	}
	queueInspector := utils.GetQueueInspector()
	queues, err := queueInspector.History("default", 30)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	maxStorageSize, err := utils.GetFileSize(utils.GetEnv("MAX_LOCALSTORAGE_SIZE"))
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	return utils.HTTPSuccessHandler(c, map[string]any{
		"version": "0.1.0",
		"total": map[string]any{
			"file_size": filesSize,
			"file_num":  len(keys),
		},
		"limit": map[string]any{
			"file_size": maxStorageSize,
		},
		"admin": map[string]any{
			"user_name": utils.GetEnv("ADMIN_NAME"),
			"email":     utils.GetEnv("ADMIN_EMAIL"),
		},
		"queue": queues,
	})
}
