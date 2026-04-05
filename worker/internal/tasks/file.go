package tasks

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"pkg/models"
	u "pkg/utils"

	"github.com/hibiken/asynq"
)

func RemoveFile(ctx context.Context, task *asynq.Task) error {
	var payload RemoveFileTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	fileInfo, err := models.GetRedisFileInfo(payload.FileId)
	if err != nil {
		return err
	}
	if fileInfo == nil {
		return nil
	}
	// 如果文件是上传文件，则需要检查是否还有分享，考虑到比如文件转换这些一次性任务产生的文件需要销毁
	if fileInfo.FileType == models.FileTypeUpload {
		shareIDs, err := models.GetRedisFileShareRelational(payload.FileId)
		if err != nil {
			return err
		}
		if len(shareIDs) > 0 {
			return nil
		}
	}

	rdb, rctx := u.GetRedisClient()
	uploadPath, err := u.GetUploadDirPath()
	if err != nil {
		return err
	}
	filePath := filepath.Join(uploadPath, payload.FileId)
	if err := rdb.Do(rctx, rdb.B().Hdel().Key("015:fileInfoMap").Field(payload.FileId).Build()).Error(); err != nil {
		return err
	}
	if err := os.RemoveAll(filePath); err != nil {
		return err
	}
	return nil
}
