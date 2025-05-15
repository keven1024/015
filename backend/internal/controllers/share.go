package controllers

import (
	"backend/internal/models"
	"backend/internal/utils"
	"backend/middleware"
	"errors"
	"time"

	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type CreateShareProps struct {
	Type models.ShareType `json:"type"`
	// ShareId string           `json:"id"`
	Config   ShareConfig `json:"config"`
	Data     string      `json:"data"`
	FileName string      `json:"file_name"`
}

type ShareConfig struct {
	ExpireAt      int      `json:"expire_time"` // 分钟
	ViewNum       int64    `json:"download_nums"`
	HasPassword   bool     `json:"has_password"`
	Password      string   `json:"password"`
	HasNotify     bool     `json:"has_notify"`
	NotifyEmail   []string `json:"notify_email"`
	HasPickupCode bool     `json:"has_pickup_code"`
}

func CreateShareInfo(c echo.Context) error {
	cc := c.(*middleware.CustomContext)

	r := new(CreateShareProps)
	if err := cc.Bind(r); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	if r.Config.ExpireAt < 1 {
		return utils.HTTPErrorHandler(c, errors.New("非法的分享过期时间"))
	}
	ExpireTime := time.Now().Add(time.Duration(r.Config.ExpireAt) * time.Minute)
	if r.Data == "" || (r.Type != models.ShareTypeFile && r.Type != models.ShareTypeText) || ExpireTime.Before(time.Now()) || r.Config.ViewNum < 1 {
		return utils.HTTPErrorHandler(c, errors.New("调用接口参数错误"))
	}

	id, err := gonanoid.New()
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	if r.Type == models.ShareTypeFile {
		fileInfo, err := models.GetRedisFileInfo(r.Data)
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		if fileInfo == nil {
			return utils.HTTPErrorHandler(c, errors.New("分享文件不存在"))
		}
		if fileInfo.FileType != models.FileTypeUpload {
			return utils.HTTPErrorHandler(c, errors.New("分享文件状态错误"))
		}
	}

	models.SetRedisShareInfo(id, models.RedisShareInfo{
		Data:        r.Data,
		Type:        r.Type,
		CreatedAt:   time.Now().Unix(),
		Owner:       cc.Auth.(string),
		ViewNum:     r.Config.ViewNum,
		Password:    r.Config.Password,
		NotifyEmail: r.Config.NotifyEmail,
		FileName:    r.FileName,
		ExpireAt:    ExpireTime.Unix(),
	})

	return utils.HTTPSuccessHandler(c, map[string]any{
		"id":            id,
		"file_name":     r.FileName,
		"download_nums": r.Config.ViewNum,
		"expire_at":     ExpireTime.Unix(),
	})
}

type GetShareProps struct {
	ShareId string `param:"id"`
}

func GetShareInfo(c echo.Context) error {
	cc := c.(*middleware.CustomContext)
	shareId := cc.Param("id")
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
		return utils.HTTPSuccessHandler(c, map[string]any{
			"id":            shareId,
			"type":          shareInfo.Type,
			"name":          shareInfo.FileName,
			"download_nums": shareInfo.ViewNum,
			"expire_at":     shareInfo.ExpireAt,
			"owner":         shareInfo.Owner,
			"size":          fileInfo.FileSize,
			"mime_type":     fileInfo.MimeType,
		})
	}

	return utils.HTTPSuccessHandler(c, map[string]any{
		"id":            shareId,
		"type":          shareInfo.Type,
		"name":          shareInfo.FileName,
		"download_nums": shareInfo.ViewNum,
		"expire_at":     shareInfo.ExpireAt,
		"owner":         shareInfo.Owner,
	})
}
