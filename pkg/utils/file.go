package utils

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
)

func GetFileId(fileHash string, fileSize int64) string {
	return fmt.Sprintf("%s_%d", fileHash, fileSize)
}

func GetFileSHA1(file io.Reader) (string, error) {
	const bufferSize = 1024 * 1000 // 1MB
	hash := sha1.New()
	buf := make([]byte, bufferSize)
	reader := bufio.NewReader(file)
	for {
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
	uploadPath := GetEnv("upload.path")
	if uploadPath == "" {
		basepath, err := os.Getwd()
		if err != nil {
			return "", err
		}
		uploadPath = filepath.Join(basepath, "uploads")
	}
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", err
	}
	return uploadPath, nil
}

func GetFileSize(size string) (uint64, error) {
	return humanize.ParseBytes(size)
}

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
