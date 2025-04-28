package utils

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func GetFileId(fileHash string, fileSize int64) string {
	return fmt.Sprintf("%s_%d", fileHash, fileSize)
}

func GetFileMd5(file *os.File) (string, error) {

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
