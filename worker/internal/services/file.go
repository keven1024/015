package services

import (
	"errors"
	"os"
	"path/filepath"
	"time"
	"worker/internal/models"
	"worker/internal/utils"
)

type GenStandardFileReturn struct {
	FileId string
	models.FileInfo
}

// 生成标准格式的file
func GenStandardFile(filePath string, mimeType string) (GenStandardFileReturn, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return GenStandardFileReturn{}, errors.New("文件不存在")
	}
	compressedFile, err := os.Open(filePath)
	if err != nil {
		return GenStandardFileReturn{}, err
	}
	defer compressedFile.Close()

	compressedFileInfo, err := compressedFile.Stat()
	if err != nil {
		return GenStandardFileReturn{}, err
	}
	compressedFileSize := compressedFileInfo.Size()

	compressedFileHash, err := utils.GetFileMd5(compressedFile)
	if err != nil {
		return GenStandardFileReturn{}, err
	}

	compressedFileId := utils.GetFileId(compressedFileHash, compressedFileSize)

	uploadPath, err := utils.GetUploadDirPath()
	if err != nil {
		return GenStandardFileReturn{}, err
	}
	newPath := filepath.Join(uploadPath, compressedFileId)
	if err := os.Rename(filePath, newPath); err != nil {
		return GenStandardFileReturn{}, err
	}
	models.SetRedisFileInfo(compressedFileId, models.RedisFileInfo{
		FileInfo: models.FileInfo{
			FileSize: compressedFileSize,
			FileHash: compressedFileHash,
			MimeType: mimeType,
		},
		FileType:  models.FileTypeUpload,
		CreatedAt: time.Now().Unix(),
	})

	return GenStandardFileReturn{
		FileId: compressedFileId,
		FileInfo: models.FileInfo{
			FileSize: compressedFileSize,
			FileHash: compressedFileHash,
			MimeType: mimeType,
		},
	}, nil
}
