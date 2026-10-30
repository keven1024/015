package utils

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/dustin/go-humanize"
)

var fileIDPattern = regexp.MustCompile(`^[0-9a-f]{40}_[0-9]+$`)

func VaildFileId(fileID string) bool {
	return fileIDPattern.MatchString(fileID)
}

func GetFileId(fileHash string, fileSize int64) (string, error) {
	fileID := fmt.Sprintf("%s_%d", fileHash, fileSize)
	if !VaildFileId(fileID) {
		return "", fmt.Errorf("invalid file id")
	}
	return fileID, nil
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
	uploadPath := GetConfig().Upload.Path
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

func GetUploadDir() (*os.Root, error) {
	path, err := GetUploadDirPath()
	if err != nil {
		return nil, err
	}
	return os.OpenRoot(path)
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
