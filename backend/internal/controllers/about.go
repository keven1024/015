package controllers

import (
	"backend/internal/utils"
	filemodel "pkg/models/file"
	u "pkg/utils"

	"github.com/labstack/echo/v5"
	"github.com/samber/lo"
)

func GetAbout(c *echo.Context) error {
	maxStorageSize, err := u.GetFileSize(u.GetConfig().Upload.Maximum)
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
		"bg_url":  u.GetConfig().About.BgURL,
		"content": u.GetEnvMap("about.content"),
		"email":   u.GetConfig().About.Email,
		"name":    u.GetConfig().About.Name,
		"url":     u.GetConfig().About.URL,
		"avatar":  u.GetConfig().About.Avatar,
		"file": map[string]any{
			"maximun": maxStorageSize,
			"current": currentFileSize,
		},
	})
}
