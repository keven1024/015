package tasks

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"worker/internal/models"
	"worker/internal/utils"

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
		rdb, ctx := utils.GetRedisClient()
		uploadPath, err := utils.GetUploadDirPath()
		if err != nil {
			return err
		}
		filePath := filepath.Join(uploadPath, payload.FileId)
		rdb.HDel(ctx, "015:fileShareRelational", payload.FileId)
		rdb.HDel(ctx, "015:fileInfoMap", payload.FileId)
		os.RemoveAll(filePath)
		return nil
	}
	models.SetRedisFileShareRelational(payload.FileId, shareIDs)
	return nil
}

func ShareNotify(ctx context.Context, task *asynq.Task) error {
	return nil
}
