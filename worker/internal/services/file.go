package services

import (
	"os"
	"path/filepath"
	"pkg/models"
	"pkg/services"
	u "pkg/utils"
	"time"
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

	fileHash, err := u.GetFileSHA1(file)
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
	redisFileInfo, err := models.SetRedisFileInfo(fileId, func(fileInfo *models.RedisFileInfo) *models.RedisFileInfo {
		fileInfo.FileInfo = models.FileInfo{
			FileSize: fileSize,
			FileHash: fileHash,
			MimeType: mimeType,
		}
		fileInfo.FileType = models.FileTypeUpload
		return fileInfo
	})
	if err != nil {
		return GenStandardFileReturn{}, err
	}
	err = services.SetFileRemoveTask(fileId, time.Duration(redisFileInfo.Expire)*time.Second)
	if err != nil {
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
