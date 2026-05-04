package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"pkg/models"
	u "pkg/utils"
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
	shareIDs, err := models.GetRedisFileShareRelational(payload.FileId)
	if err != nil {
		return err
	}
	shareIDs = lo.Filter(shareIDs, func(x string, _ int) bool {
		return x != payload.ShareId
	})
	if len(shareIDs) == 0 {
		rdb, ctx := u.GetRedisClient()
		uploadPath, err := u.GetUploadDirPath()
		if err != nil {
			return err
		}
		filePath := filepath.Join(uploadPath, payload.FileId)
		if err := rdb.Do(ctx, rdb.B().Hdel().Key("015:fileShareRelational").Field(payload.FileId).Build()).Error(); err != nil {
			return err
		}
		if err := rdb.Do(ctx, rdb.B().Hdel().Key("015:fileInfoMap").Field(payload.FileId).Build()).Error(); err != nil {
			return err
		}
		if err := os.RemoveAll(filePath); err != nil {
			return err
		}
		return nil
	}
	if err := models.SetRedisFileShareRelational(payload.FileId, shareIDs); err != nil {
		return err
	}
	return nil
}

func ShareNotify(ctx context.Context, task *asynq.Task) error {
	var payload ShareNotifyTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	shareInfo, err := models.GetRedisShareInfo(payload.ShareId)
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

	for _, email := range shareInfo.NotifyEmails {
		if err := services.SendEmail(email, services.EmailTemplateData{
			Locale:    shareInfo.Locale,
			ShareType: shareInfo.Type,
			FileName:  shareInfo.FileName,
			IP:        payload.IP,
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
