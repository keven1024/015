package services

import (
	"errors"
	"os"
	"path/filepath"
	"worker/internal/models"
	"worker/internal/utils"
)

// 生成标准格式的file
func GenStandardFile(filePath string, mimeType string) (string, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", errors.New("文件不存在")
	}
	compressedFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer compressedFile.Close()

	compressedFileInfo, err := compressedFile.Stat()
	if err != nil {
		return "", err
	}
	compressedFileSize := compressedFileInfo.Size()

	compressedFileHash, err := utils.GetFileMd5(compressedFile)
	if err != nil {
		return "", err
	}

	compressedFileId := utils.GetFileId(compressedFileHash, compressedFileSize)

	uploadPath, err := utils.GetUploadDirPath()
	if err != nil {
		return "", err
	}
	newPath := filepath.Join(uploadPath, compressedFileId)
	if err := os.Rename(filePath, newPath); err != nil {
		return "", err
	}
	models.SetRedisFileInfo(compressedFileId, models.RedisFileInfo{
		FileInfo: models.FileInfo{
			FileSize: compressedFileSize,
			FileHash: compressedFileHash,
			MimeType: mimeType,
		},
	})

	return compressedFileId, nil
}
