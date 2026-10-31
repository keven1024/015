package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"pkg/geoip"
	relationmodel "pkg/models/file_share_relational"
	sharemodel "pkg/models/share"
	pkgservices "pkg/services"
	"worker/internal/services"

	"github.com/hibiken/asynq"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type ShareRemoveTaskPayload struct {
	ShareId string   `json:"share_id"`
	FileIds []string `json:"file_ids"`
}

func RemoveShare(ctx context.Context, task *asynq.Task) error {
	var payload ShareRemoveTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	if ok := lo.EveryBy(payload.FileIds, func(fileId string) bool {
		shareIDs, err := relationmodel.SetRedisFileShareRelational(fileId, func(shareIDs []string) []string {
			return lo.Filter(shareIDs, func(x string, _ int) bool {
				return x != payload.ShareId
			})
		})
		if err != nil {
			zap.L().Error("remove share relation failed",
				zap.String("file_id", fileId),
				zap.String("share_id", payload.ShareId),
				zap.Error(err),
			)
			return false
		}
		if len(shareIDs) == 0 {
			if err := pkgservices.SetFileRemoveTask(fileId, 0); err != nil {
				zap.L().Error("enqueue file removal task failed",
					zap.String("file_id", fileId),
					zap.String("share_id", payload.ShareId),
					zap.Error(err),
				)
				return false
			}
		}
		return true
	}); !ok {
		return fmt.Errorf("remove share files failed for share %s", payload.ShareId)
	}
	return nil
}

// 通知策略：Webhook 和 Email 分别尝试，任一渠道至少有一个目标发送成功，
// 就认为本次通知成功；只有两个渠道全部发送失败时才返回错误并触发重试。
func ShareNotify(ctx context.Context, task *asynq.Task) error {
	var payload ShareNotifyTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	shareInfo, err := sharemodel.GetRedisShareInfo(payload.ShareId)
	if err != nil || shareInfo == nil {
		return err
	}

	webhookOK := lo.Reduce(shareInfo.NotifyWebhooks, func(ok bool, webhook sharemodel.NotifyWebhook, _ int) bool {
		if err := services.SendWebhook(webhook); err != nil {
			zap.L().Error("send share notification webhook failed",
				zap.String("share_id", payload.ShareId),
				zap.String("url", webhook.URL),
				zap.Error(err),
			)
			return ok
		}
		return true
	}, false)

	region := "-"
	if info := geoip.GetIpGeoInfo(payload.IP); info != nil {
		region = info.Emoji + " " + info.Country.Country.Names.English
	}

	displayName := lo.Substring(shareInfo.Text, 0, 7) + "..."
	if shareInfo.Type == sharemodel.ShareTypeFile && len(shareInfo.Files) > 0 {
		displayName = shareInfo.Files[0].FileName
	}

	emailOK := lo.Reduce(shareInfo.NotifyEmails, func(ok bool, email string, _ int) bool {
		if err := services.SendEmail(email, services.EmailTemplateData{
			Locale:    shareInfo.Locale,
			FileName:  displayName,
			IP:        payload.IP,
			Region:    region,
			ShareType: shareInfo.Type,
		}); err != nil {
			zap.L().Error("send share notification email failed",
				zap.String("share_id", payload.ShareId),
				zap.String("email", email),
				zap.Error(err),
			)
			return ok
		}
		return true
	}, false)

	// 只有所有通知目标都失败时才让 Asynq 重试；部分成功视为成功。
	if !webhookOK && !emailOK {
		return fmt.Errorf("share notification failed for share %s", payload.ShareId)
	}
	return nil
}
