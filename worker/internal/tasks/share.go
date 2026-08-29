package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"pkg/geoip"
	relationmodel "pkg/models/file_share_relational"
	sharemodel "pkg/models/share"
	pkgservices "pkg/services"
	"worker/internal/services"

	"github.com/hibiken/asynq"
	"github.com/samber/lo"
)

type ShareRemoveTaskPayload struct {
	ShareId string `json:"share_id"`
	FileId  string `json:"file_id"`
}

func RemoveShare(ctx context.Context, task *asynq.Task) error {
	var payload ShareRemoveTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	shareIDs, err := relationmodel.SetRedisFileShareRelational(payload.FileId, func(shareIDs []string) []string {
		return lo.Filter(shareIDs, func(x string, _ int) bool {
			return x != payload.ShareId
		})
	})
	if err != nil {
		return err
	}
	if len(shareIDs) == 0 {
		return pkgservices.SetFileRemoveTask(payload.FileId, 0)
	}
	return nil
}

func ShareNotify(ctx context.Context, task *asynq.Task) error {
	var payload ShareNotifyTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	shareInfo, err := sharemodel.GetRedisShareInfo(payload.ShareId)
	if err != nil || shareInfo == nil {
		return err
	}

	var errs []error
	successCount := 0

	for _, webhook := range shareInfo.NotifyWebhooks {
		if err := services.SendWebhook(webhook); err != nil {
			errs = append(errs, err)
			continue
		}
		successCount++
	}

	region := "-"
	if info := geoip.GetIpGeoInfo(payload.IP); info != nil {
		region = info.Emoji + " " + info.Country.Country.Names.English
	}

	displayName := lo.Substring(shareInfo.Text, 0, 7) + "..."
	if shareInfo.Type == sharemodel.ShareTypeFile && len(shareInfo.Files) > 0 {
		displayName = shareInfo.Files[0].FileName
	}

	for _, email := range shareInfo.NotifyEmails {
		if err := services.SendEmail(email, services.EmailTemplateData{
			Locale:    shareInfo.Locale,
			FileName:  displayName,
			IP:        payload.IP,
			Region:    region,
			ShareType: shareInfo.Type,
		}); err != nil {
			errs = append(errs, err)
			continue
		}
		successCount++
	}

	if successCount > 0 || len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("all share notify targets failed: %w", errors.Join(errs...))
}
