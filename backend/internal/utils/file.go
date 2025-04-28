package utils

import "fmt"

func GetFileId(fileHash string, fileSize int64) string {
	return fmt.Sprintf("%s_%d", fileHash, fileSize)
}
