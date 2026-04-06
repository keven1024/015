package services

import (
	"os"
	"path/filepath"
	"pkg/models"
	"pkg/services"
	u "pkg/utils"
	"time"

	"github.com/spf13/cast"
)

type GenStandardFileReturn struct {
	FileId string
	models.FileInfo
}

// 生成标准格式的file
func GenStandardFile(filePath string, mimeType string) (GenStandardFileReturn, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return GenStandardFileReturn{}, ErrFileNotFound
	}
	file, err := os.Open(filePath)
	if err != nil {
		return GenStandardFileReturn{}, err
	}
	defer file.Close() //nolint:errcheck

	fileInfo, err := file.Stat()
	if err != nil {
		return GenStandardFileReturn{}, err
	}
	fileSize := fileInfo.Size()

	fileHash, err := u.GetFileMd5(file)
	if err != nil {
		return GenStandardFileReturn{}, err
	}

	fileId := u.GetFileId(fileHash, fileSize)

	uploadPath, err := u.GetUploadDirPath()
	if err != nil {
		return GenStandardFileReturn{}, err
	}
	newPath := filepath.Join(uploadPath, fileId)
	if err := os.Rename(filePath, newPath); err != nil {
		return GenStandardFileReturn{}, err
	}
	expire := cast.ToInt64(u.GetEnvWithDefault("upload.remove_expire", "2")) * 3600
	err = services.SetFileRemoveTask(fileId, time.Duration(expire)*time.Second)
	if err != nil {
		return GenStandardFileReturn{}, err
	}
	if err := models.SetRedisFileInfo(fileId, func(fileInfo *models.RedisFileInfo) *models.RedisFileInfo {
		fileInfo.FileInfo = models.FileInfo{
			FileSize: fileSize,
			FileHash: fileHash,
			MimeType: mimeType,
		},
		FileType:  models.FileTypeUpload,
		CreatedAt: time.Now().Unix(),
		Expire:    expire,
	}); err != nil {
		return GenStandardFileReturn{}, err
	}
	return GenStandardFileReturn{
		FileId: fileId,
		FileInfo: models.FileInfo{
			FileSize: fileSize,
			FileHash: fileHash,
			MimeType: mimeType,
		},
	}, nil
}
