package tasks

import (
	"context"
	"encoding/json"
	"math"
	filemodel "pkg/models/file"
	pkgservices "pkg/services"
	u "pkg/utils"
	"time"

	"github.com/hibiken/asynq"
)

const uploadRemoveRetryDelay = 30 * time.Minute

// RemoveUpload 删除废弃的临时上传切片。
// 当切片进度领先于其时间进度的一半时，它会保持初始化上传，
// 30 分钟后再次检查。否则它将实际的清理工作委托给file:remove
func RemoveUpload(_ context.Context, task *asynq.Task) error {
	var payload RemoveUploadTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	fileInfo, err := filemodel.GetRedisFileInfo(payload.FileId)
	if err != nil {
		return err
	}
	if fileInfo == nil || fileInfo.FileType != filemodel.FileTypeInit {
		return nil
	}
	if fileInfo.Expire <= 0 || fileInfo.ChunkSize <= 0 || fileInfo.FileSize <= 0 {
		return pkgservices.SetFileRemoveTask(payload.FileId, 0)
	}

	root, err := u.GetUploadDir()
	if err != nil {
		return err
	}
	defer root.Close() //nolint:errcheck

	sliceList, err := pkgservices.GetFileSliceList(payload.FileId, root)
	if err != nil {
		return err
	}

	totalSlices := int(math.Ceil(float64(fileInfo.FileSize) / float64(fileInfo.ChunkSize)))
	if totalSlices == 0 {
		return pkgservices.SetFileRemoveTask(payload.FileId, 0)
	}

	elapsed := time.Now().Unix() - fileInfo.CreatedAt
	timeProgress := float64(elapsed) / float64(fileInfo.Expire)
	uploadProgress := float64(len(sliceList)) / float64(totalSlices)
	if uploadProgress > timeProgress/2 {
		return pkgservices.SetUploadRemoveTask(payload.FileId, uploadRemoveRetryDelay)
	}

	return pkgservices.SetFileRemoveTask(payload.FileId, 0)
}
