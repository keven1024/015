package tasks

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	filemodel "pkg/models/file"
	relationmodel "pkg/models/file_share_relational"
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
	fileInfo, err := filemodel.GetRedisFileInfo(payload.FileId)
	if err != nil {
		return err
	}
	if fileInfo == nil {
		return nil
	}
	// 如果文件是上传文件，则需要检查是否还有分享，考虑到比如文件转换这些一次性任务产生的文件需要销毁
	if fileInfo.FileType == filemodel.FileTypeUpload {
		shareIDs, err := relationmodel.GetRedisFileShareRelational(payload.FileId)
		if err != nil {
			return err
		}
		if len(shareIDs) > 0 {
			return nil
		}
	}

	uploadPath, err := u.GetUploadDirPath()
	if err != nil {
		return err
	}
	filePath := filepath.Join(uploadPath, payload.FileId)
	// 如果是临时文件删除文件夹
	if fileInfo.FileType == filemodel.FileTypeInit {
		filePath += "_tmp"
	}
	if err := filemodel.DelRedisFileInfo(payload.FileId); err != nil {
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

	allFileInfo, err := filemodel.GetRedisFileInfoAll()
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
		info, err := filemodel.JsonFileInfoToDomain(rawInfo)
		if err != nil {
			continue
		}

		// Case 2: init 狀態且已過期
		if info.FileType == filemodel.FileTypeInit && info.CreatedAt+info.Expire < now {
			if err := pkgservices.SetFileRemoveTask(fileId, 0); err != nil {
				return err
			}
			continue
		}

		// Case 3: 已完成上傳但無 share 關係
		if info.FileType == filemodel.FileTypeUpload {
			shareIDs, err := relationmodel.GetRedisFileShareRelational(fileId)
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
