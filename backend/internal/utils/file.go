package utils

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"

	goUnits "github.com/docker/go-units"
)

func GetFileId(fileHash string, fileSize int64) string {
	return fmt.Sprintf("%s_%d", fileHash, fileSize)
}

func GetFileMd5(file io.Reader) (string, error) {

	const bufferSize = 1024 * 1000 // 1MB

	hash := md5.New()
	for buf, reader := make([]byte, bufferSize), bufio.NewReader(file); ; {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		hash.Write(buf[:n])
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func GetUploadDirPath() (string, error) {
	basepath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	finalPath := filepath.Join(basepath, "uploads")
	uploadPath := GetEnvWithDefault("UPLOAD_PATH", finalPath)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", err
	}
	return uploadPath, nil
}

func GetFileSize(size string) (int64, error) {
	s, err := goUnits.FromHumanSize(size)
	return s, err
}
