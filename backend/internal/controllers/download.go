package controllers

import (
	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/utils"
	"backend/middleware"
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

func DownloadShare(c echo.Context) error {
	cc := c.(*middleware.CustomContext)
	shareId := cc.Param("id")
	password := cc.QueryParam("password")

	if shareId == "" {
		return utils.HTTPErrorHandler(c, errors.New("缺少分享ID"))
	}

	shareInfo, err := models.GetRedisShareInfo(shareId)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	if shareInfo == nil {
		return utils.HTTPErrorHandler(c, errors.New("分享不存在"))
	}
	if shareInfo.Password != "" && shareInfo.Password != password {
		return utils.HTTPErrorHandler(c, errors.New("分享密码错误"))
	}

	if shareInfo.Type == models.ShareTypeFile {
		fileInfo, err := models.GetRedisFileInfo(shareInfo.Data)
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		if fileInfo == nil {
			return utils.HTTPErrorHandler(c, errors.New("分享文件不存在"))
		}
		if fileInfo.FileType != models.FileTypeUpload {
			return utils.HTTPErrorHandler(c, errors.New("分享文件状态错误"))
		}
		uploadPath, err := services.GetUploadDirPath()
		if err != nil {
			return err
		}
		return cc.Attachment(fmt.Sprintf("%s/%s", uploadPath, utils.GetFileId(fileInfo.FileHash, fileInfo.FileSize)), shareInfo.FileName)
	}

	return utils.HTTPSuccessHandler(c, map[string]any{
		"data": shareInfo.Data,
	})
}
