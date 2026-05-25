package tasks

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"pkg/models"
	pkgservices "pkg/services"
	u "pkg/utils"
	"strings"
	"time"

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

	rdb := u.GetRedisClient()
	uploadPath, err := u.GetUploadDirPath()
	if err != nil {
		return err
	}
	filePath := filepath.Join(uploadPath, payload.FileId)
	// 如果是临时文件删除文件夹
	if fileInfo.FileType == models.FileTypeInit {
		filePath += "_tmp"
	}
	if err := rdb.Do(ctx, rdb.B().Hdel().Key("015:fileInfoMap").Field(payload.FileId).Build()).Error(); err != nil {
		return err
	}
	if err := os.RemoveAll(filePath); err != nil {
		return err
	}
	return nil
}

func FileJanitor(_ context.Context, _ *asynq.Task) error {
	uploadPath, err := u.GetUploadDirPath()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(uploadPath)
	if err != nil {
		return err
	}

	allFileInfo, err := models.GetRedisFileInfoAll()
	if err != nil {
		return err
	}

	// Case 1: 本地有但 fileInfoMap 無 → 直接刪除
	for _, entry := range entries {
		name := entry.Name()
		fileId := strings.TrimSuffix(name, "_tmp")
		if _, exists := allFileInfo[fileId]; !exists {
			if err := os.RemoveAll(filepath.Join(uploadPath, name)); err != nil {
				return err
			}
		}
	}

	// Case 2 & 3: 遍歷 fileInfoMap
	now := time.Now().Unix()
	for fileId, rawInfo := range allFileInfo {
		var info models.RedisFileInfo
		if err := json.Unmarshal([]byte(rawInfo), &info); err != nil {
			continue
		}

		// Case 2: init 狀態且已過期
		if info.FileType == models.FileTypeInit && info.CreatedAt+info.Expire < now {
			if err := pkgservices.SetFileRemoveTask(fileId, 0); err != nil {
				return err
			}
			continue
		}

		// Case 3: 已完成上傳但無 share 關係
		if info.FileType == models.FileTypeUpload {
			shareIDs, err := models.GetRedisFileShareRelational(fileId)
			if err != nil {
				return err
			}
			if len(shareIDs) == 0 {
				if err := pkgservices.SetFileRemoveTask(fileId, 0); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
