package controllers

import (
	"backend/internal/models"
	"backend/internal/utils"
	"backend/middleware"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/spf13/cast"
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
	password := ""
	if r.Config.Password != "" {
		hash, err := utils.GeneratePasswordHash(r.Config.Password)
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		password = hash
	}

	models.SetRedisShareInfo(id, models.RedisShareInfo{
		Data:      r.Data,
		Type:      r.Type,
		CreatedAt: time.Now().Unix(),
		Owner:     cc.Auth.(string),
		ViewNum:   r.Config.ViewNum,
		Password:  password,
		// NotifyEmail: r.Config.NotifyEmail,
		FileName: r.FileName,
		ExpireAt: ExpireTime.Unix(),
	})
	var pickupCode string
	if r.Config.HasPickupCode {
		for {
			pickupCode = utils.GeneratePickupCode()
			ok, err := models.SetRedisPickupData(pickupCode, id)
			if err != nil {
				return utils.HTTPErrorHandler(c, err)
			}
			if !ok {
				continue
			}
			break
		}
	}

	if r.Type == models.ShareTypeFile {
		shareIDs, err := models.GetRedisFileShareRelational(r.Data)
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		shareIDs = append(shareIDs, id)
		models.SetRedisFileShareRelational(r.Data, shareIDs)
		client := utils.GetQueueClient()
		json, err := json.Marshal(map[string]any{"share_id": id, "file_id": r.Data})
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		// 这里延时分享过期时间基础上加下载窗口期后1小时删除，防止用户过期前几分钟才开始下载，下载一半文件不见了
		downloadWindow := utils.GetEnvWithDefault("share.download_window", "12")
		deleteTime := time.Duration(r.Config.ExpireAt)*time.Minute + cast.ToDuration(downloadWindow+"h") + 1*time.Hour
		_, err = client.Enqueue(asynq.NewTask("share:remove", json), asynq.ProcessIn(deleteTime))
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
	}

	// 统计分享数
	currentDate := time.Now().Format("2006-01-02")
	statData, _ := models.GetRedisStat(currentDate)
	if statData == nil {
		statData = &models.StatData{
			FileSize:    0,
			FileNum:     0,
			ShareNum:    0,
			DownloadNum: 0,
		}
	}
	statData.ShareNum += 1
	models.SetRedisStat(currentDate, *statData)

	return utils.HTTPSuccessHandler(c, map[string]any{
		"id":            id,
		"file_name":     r.FileName,
		"download_nums": r.Config.ViewNum,
		"expire_at":     ExpireTime.Unix(),
		"pickup_code":   pickupCode,
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
	if shareInfo == nil || shareInfo.ViewNum < 1 {
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
			"has_password":  shareInfo.Password != "",
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
		"has_password":  shareInfo.Password != "",
		"expire_at":     shareInfo.ExpireAt,
		"owner":         shareInfo.Owner,
	})
}

func GetShareByPickupCode(c echo.Context) error {
	cc := c.(*middleware.CustomContext)
	pickupCode := cc.Param("code")
	if pickupCode == "" {
		return utils.HTTPErrorHandler(c, errors.New("缺少提取码"))
	}
	shareId, err := models.GetRedisPickupData(strings.ToUpper(pickupCode))
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	if shareId == "" {
		return utils.HTTPErrorHandler(c, errors.New("分享不存在"))
	}
	return utils.HTTPSuccessHandler(c, map[string]any{
		"share_id": shareId,
	})
}
