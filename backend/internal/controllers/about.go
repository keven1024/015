package controllers

import (
	"backend/internal/utils"
	filemodel "pkg/models/file"
	u "pkg/utils"

	"github.com/labstack/echo/v5"
	"github.com/samber/lo"
)

func GetAbout(c *echo.Context) error {
	cfg := u.GetConfig()
	maxStorageSize, err := u.GetFileSize(cfg.Upload.Maximum)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	fileInfoMap, err := filemodel.GetRedisFileInfoAll()
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	currentFileSize := lo.Reduce(lo.Values(fileInfoMap), func(agg int64, item string, _ int) int64 {
		fileInfo, err := filemodel.JsonFileInfoToDomain(item)
		if err != nil {
			return agg
		}
		return agg + fileInfo.FileSize
	}, 0)

	return utils.HTTPSuccessHandler(c, map[string]any{
		"bg_url":  cfg.About.BgURL,
		"content": cfg.About.Content,
		"email":   cfg.About.Email,
		"name":    cfg.About.Name,
		"url":     cfg.About.URL,
		"avatar":  cfg.About.Avatar,
		"file": map[string]any{
			"maximun": maxStorageSize,
			"current": currentFileSize,
		},
	})
}
