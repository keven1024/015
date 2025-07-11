package services

import (
	"backend/internal/utils"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func CreateFileSlice(fileSlice io.Reader, fileId string, fileIndex int64) error {
	uploadPath, err := utils.GetUploadDirPath()
	if err != nil {
		return err
	}

	filePath := filepath.Join(uploadPath, fmt.Sprintf("%s_%s", fileId, "tmp"))
	if err := os.MkdirAll(filePath, 0755); err != nil {
		return err
	}

	dst, err := os.Create(filepath.Join(filePath, fmt.Sprintf("%d", fileIndex)))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, fileSlice); err != nil {
		return err
	}

	return nil
}

// MergeFileSlices 合并文件切片
func MergeFileSlices(slicesPath string, mergeFilePath string) error {
	// 创建最终文件
	destFile, err := os.Create(mergeFilePath)
	if err != nil {
		return fmt.Errorf("创建合并文件失败: %v", err)
	}
	defer destFile.Close()

	// 读取目录下的所有文件
	files, err := os.ReadDir(slicesPath)
	if err != nil {
		return fmt.Errorf("读取切片目录失败: %v", err)
	}

	// 按照索引排序文件切片
	sliceFiles := make([]string, len(files))
	for _, file := range files {
		index, err := strconv.Atoi(file.Name())
		if err != nil {
			return fmt.Errorf("无效的切片文件名: %v", err)
		}
		sliceFiles[index-1] = filepath.Join(slicesPath, file.Name())
	}

	// 合并文件
	buffer := make([]byte, 4*1024*1024) // 4MB buffer
	for _, sliceFile := range sliceFiles {
		sf, err := os.Open(sliceFile)
		if err != nil {
			return fmt.Errorf("打开切片文件失败: %v", err)
		}

		for {
			n, err := sf.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				sf.Close()
				return fmt.Errorf("读取切片文件失败: %v", err)
			}

			if _, err := destFile.Write(buffer[:n]); err != nil {
				sf.Close()
				return fmt.Errorf("写入合并文件失败: %v", err)
			}
		}

		sf.Close()
	}
	defer os.RemoveAll(slicesPath)
	return nil
}
