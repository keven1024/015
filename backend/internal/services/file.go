package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func CreateFileSlice(fileId string, uploadPath string, fileSlice io.Reader, fileIndex int64) (string, error) {
	filePath := filepath.Join(uploadPath, fmt.Sprintf("%s_%s", fileId, "tmp"))
	if err := os.MkdirAll(filePath, 0755); err != nil {
		return "", err
	}

	dst, err := os.Create(filepath.Join(filePath, fmt.Sprintf("%d", fileIndex)))
	if err != nil {
		return "", err
	}
	defer dst.Close() //nolint:errcheck

	if _, err = io.Copy(dst, fileSlice); err != nil {
		return "", err
	}

	return filePath, nil
}

func GetFileSliceList(fileId string, uploadPath string) ([]int, error) {
	slicesPath := filepath.Join(uploadPath, fmt.Sprintf("%s_%s", fileId, "tmp"))
	files, err := os.ReadDir(slicesPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("读取切片目录失败: %v", err)
	}
	fileSliceList := []int{}
	for _, file := range files {
		index, err := strconv.Atoi(file.Name())
		if err != nil {
			return nil, fmt.Errorf("无效的切片文件名: %v", err)
		}
		fileSliceList = append(fileSliceList, index)
	}
	sort.Ints(fileSliceList)
	return fileSliceList, nil
}

// MergeFileSlices 合并文件切片
func MergeFileSlices(fileId string, uploadPath string) (string, error) {
	mergeFilePath := filepath.Join(uploadPath, fileId)
	slicesPath := filepath.Join(uploadPath, fmt.Sprintf("%s_%s", fileId, "tmp"))
	defer os.RemoveAll(slicesPath) //nolint:errcheck
	// 创建最终文件
	destFile, err := os.Create(mergeFilePath)
	if err != nil {
		return "", fmt.Errorf("创建合并文件失败: %v", err)
	}
	defer destFile.Close() //nolint:errcheck

	fileSliceList, err := GetFileSliceList(fileId, uploadPath)
	if err != nil {
		return "", err
	}

	// 合并文件
	buffer := make([]byte, 4*1024*1024) // 4MB buffer
	for _, index := range fileSliceList {
		sliceFilePath := filepath.Join(slicesPath, fmt.Sprintf("%d", index))
		sf, err := os.Open(sliceFilePath)
		if err != nil {
			return "", fmt.Errorf("打开切片文件失败: %v", err)
		}

		if _, err := io.CopyBuffer(destFile, sf, buffer); err != nil {
			sf.Close() //nolint:errcheck
			return "", fmt.Errorf("合并切片文件失败: %v", err)
		}
		if err := sf.Close(); err != nil {
			return "", fmt.Errorf("关闭切片文件失败: %v", err)
		}
	}
	return mergeFilePath, nil
}
