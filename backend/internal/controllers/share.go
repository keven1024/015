package controllers

import (
	"backend/internal/utils"
	"encoding/json"
	filemodel "pkg/models/file"
	relationmodel "pkg/models/file_share_relational"
	pickupcodemodel "pkg/models/pickupcode"
	sharemodel "pkg/models/share"
	statmodel "pkg/models/stat"
	u "pkg/utils"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type CreateShareProps struct {
	Type sharemodel.ShareType `json:"type"`
	// ShareId string           `json:"id"`
	Config ShareConfig                `json:"config"`
	Text   string                     `json:"text"`
	Files  []sharemodel.ShareFileData `json:"files"`
}

type ShareConfig struct {
	ExpireAt       int                        `json:"expire_time"` // 分钟
	ViewNum        int64                      `json:"download_nums"`
	HasPassword    bool                       `json:"has_password"`
	Password       string                     `json:"password"`
	HasNotify      bool                       `json:"has_notify"`
	NotifyTypes    []string                   `json:"notify_types"`
	NotifyEmails   []string                   `json:"notify_emails"`
	NotifyWebhooks []sharemodel.NotifyWebhook `json:"notify_webhooks"`
	Locale         string                     `json:"locale"`
	HasPickupCode  bool                       `json:"has_pickup_code"`
}

func CreateShareInfo(c *echo.Context) error {
	owner, _ := echo.ContextGet[string](c, "auth")

	r := new(CreateShareProps)
	if err := c.Bind(r); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	if r.Config.ExpireAt < 1 {
		return utils.HTTPErrorHandler(c, ErrInvalidRequest)
	}
	ExpireTime := time.Now().Add(time.Duration(r.Config.ExpireAt) * time.Minute)
	if (r.Type != sharemodel.ShareTypeFile && r.Type != sharemodel.ShareTypeText) || ExpireTime.Before(time.Now()) || r.Config.ViewNum < 1 {
		return utils.HTTPErrorHandler(c, ErrInvalidRequest)
	}

	id, err := gonanoid.New()
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	if r.Type == sharemodel.ShareTypeFile {
		if len(r.Files) == 0 {
			return utils.HTTPErrorHandler(c, ErrInvalidRequest)
		}
		var fileInfoErr error
		if !lo.EveryBy(r.Files, func(file sharemodel.ShareFileData) bool {
			if file.Id == "" {
				return false
			}
			fileInfo, err := filemodel.GetRedisFileInfo(file.Id)
			if err != nil {
				fileInfoErr = err
				return false
			}
			if fileInfo == nil {
				return false
			}
			if fileInfo.FileType != filemodel.FileTypeUpload {
				return false
			}
			return true
		}) {
			if fileInfoErr != nil {
				return utils.HTTPErrorHandler(c, fileInfoErr)
			}
			return utils.HTTPErrorHandler(c, ErrInvalidShareFileData)
		}
	} else {
		if r.Text == "" {
			return utils.HTTPErrorHandler(c, ErrInvalidRequest)
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

	var notifyEmails []string
	var notifyWebhooks []sharemodel.NotifyWebhook
	if r.Config.HasNotify {
		if !lo.EveryBy(r.Config.NotifyTypes, func(nt string) bool {
			return nt == "email" || nt == "webhook"
		}) {
			return utils.HTTPErrorHandler(c, ErrInvalidRequest)
		}
		if lo.Contains(r.Config.NotifyTypes, "email") {
			notifyEmails = r.Config.NotifyEmails
		}
		if lo.Contains(r.Config.NotifyTypes, "webhook") {
			notifyWebhooks = r.Config.NotifyWebhooks
		}
	}

	_, err = sharemodel.SetRedisShareInfo(id, func(shareInfo *sharemodel.RedisShareInfo) *sharemodel.RedisShareInfo {
		shareInfo.Text = r.Text
		shareInfo.Files = r.Files
		shareInfo.Type = r.Type
		shareInfo.CreatedAt = time.Now().Unix()
		shareInfo.Owner = owner
		shareInfo.ViewNum = r.Config.ViewNum
		shareInfo.Password = password
		shareInfo.NotifyEmails = notifyEmails
		shareInfo.NotifyWebhooks = notifyWebhooks
		shareInfo.Locale = r.Config.Locale
		shareInfo.ExpireAt = ExpireTime.Unix()
		return shareInfo
	})
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	var pickupCode string
	if r.Config.HasPickupCode {
		pickupCodeExpireAt := time.Now().Add(24 * time.Hour).Unix()
		for {
			pickupCode = utils.GeneratePickupCode()
			ok, err := pickupcodemodel.SetRedisPickupData(pickupCode, id)
			if err != nil {
				return utils.HTTPErrorHandler(c, err)
			}
			if !ok {
				continue
			}
			break
		}
		_, err = sharemodel.SetRedisShareInfo(id, func(shareInfo *sharemodel.RedisShareInfo) *sharemodel.RedisShareInfo {
			shareInfo.PickupCode = pickupCode
			shareInfo.PickupCodeExpireAt = pickupCodeExpireAt
			return shareInfo
		})
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
	}

	if r.Type == sharemodel.ShareTypeFile {
		fileIDs := lo.Map(r.Files, func(file sharemodel.ShareFileData, _ int) string {
			return file.Id
		})
		for _, file := range r.Files {
			shareIDs, err := relationmodel.GetRedisFileShareRelational(file.Id)
			if err != nil {
				return utils.HTTPErrorHandler(c, err)
			}
			shareIDs = lo.Uniq(lo.Concat(shareIDs, []string{id}))
			err = relationmodel.SetRedisFileShareRelational(file.Id, shareIDs)
			if err != nil {
				return utils.HTTPErrorHandler(c, err)
			}
		}
		client := u.GetQueueClient()
		payload, err := json.Marshal(map[string]any{"share_id": id, "file_ids": fileIDs})
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		// 这里延时分享过期时间基础上加下载窗口期后1小时删除，防止用户过期前几分钟才开始下载，下载一半文件不见了
		downloadWindow := u.GetEnvWithDefault("share.download_window", "12")
		deleteTime := time.Duration(r.Config.ExpireAt)*time.Minute + cast.ToDuration(downloadWindow+"h") + 1*time.Hour
		_, err = client.Enqueue(asynq.NewTask("share:remove", payload), asynq.ProcessIn(deleteTime))
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
	}

	// 统计分享数
	currentDate := time.Now().Format("2006-01-02")
	_, err = statmodel.SetRedisStat(currentDate, func(stat *statmodel.StatData) *statmodel.StatData {
		stat.ShareNum += 1
		return stat
	})
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	return utils.HTTPSuccessHandler(c, map[string]any{
		"id":            id,
		"files":         r.Files,
		"download_nums": r.Config.ViewNum,
		"expire_at":     ExpireTime.Unix(),
		"pickup_code":   pickupCode,
	})
}

type GetShareProps struct {
	ShareId string `param:"id"`
}

func GetShareInfo(c *echo.Context) error {
	shareId := c.Param("id")
	if shareId == "" {
		return utils.HTTPErrorHandler(c, ErrInvalidRequest)
	}

	shareInfo, err := sharemodel.GetRedisShareInfo(shareId)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	if shareInfo == nil || shareInfo.ExpireAt <= time.Now().Unix() || shareInfo.ViewNum < 1 {
		return utils.HTTPErrorHandler(c, ErrShareNotFound)
	}
	owner, _ := echo.ContextGet[string](c, "auth")
	isOwner := owner != "" && owner == shareInfo.Owner
	response := map[string]any{
		"id":            shareId,
		"type":          shareInfo.Type,
		"download_nums": shareInfo.ViewNum,
		"has_password":  shareInfo.Password != "",
		"has_notify":    len(shareInfo.NotifyEmails) > 0 || len(shareInfo.NotifyWebhooks) > 0,
		"expire_at":     shareInfo.ExpireAt,
		"owner":         shareInfo.Owner,
		"is_owner":      isOwner,
	}
	if isOwner {
		response["pickup_code"] = shareInfo.PickupCode
		response["pickup_code_expire_at"] = shareInfo.PickupCodeExpireAt
		response["notify_emails"] = shareInfo.NotifyEmails
		response["notify_webhooks"] = shareInfo.NotifyWebhooks
	}

	if shareInfo.Type == sharemodel.ShareTypeFile {
		shareFiles := shareInfo.Files
		files := make([]map[string]any, 0, len(shareFiles))
		for _, file := range shareFiles {
			fileInfo, err := filemodel.GetRedisFileInfo(file.Id)
			if err != nil {
				return utils.HTTPErrorHandler(c, err)
			}
			if fileInfo == nil {
				return utils.HTTPErrorHandler(c, ErrShareFileNotFound)
			}
			if fileInfo.FileType != filemodel.FileTypeUpload {
				return utils.HTTPErrorHandler(c, ErrInvalidShareFileState)
			}
			files = append(files, map[string]any{
				"id":        file.Id,
				"file_name": file.FileName,
				"size":      fileInfo.FileSize,
				"mime_type": fileInfo.MimeType,
			})
		}
		response["files"] = files
		return utils.HTTPSuccessHandler(c, response)
	}
	if isOwner {
		response["text"] = shareInfo.Text
	}

	return utils.HTTPSuccessHandler(c, response)
}

func GetShareByPickupCode(c *echo.Context) error {
	pickupCode := c.Param("code")
	if pickupCode == "" {
		return utils.HTTPErrorHandler(c, ErrInvalidRequest)
	}
	shareId, err := pickupcodemodel.GetRedisPickupData(strings.ToUpper(pickupCode))
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	if shareId == "" {
		return utils.HTTPErrorHandler(c, ErrShareNotFound)
	}
	return utils.HTTPSuccessHandler(c, map[string]any{
		"share_id": shareId,
	})
}
