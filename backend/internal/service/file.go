package service

import (
	"backend/internal/models"
	"backend/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)

func GetFileInfo() {

}

func GetUploadDirPath() (string, error) {
	basepath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	finalPath := filepath.Join(basepath, "upload")
	if err := os.MkdirAll(finalPath, 0755); err != nil {
		return "", err
	}
	return finalPath, nil
}

func GetRedisFileInfo(fileId string) (models.RedisFileInfo, error) {
	rdb, ctx := utils.GetRedisClient()
	fileInfoUnmarshalData, _ := rdb.HGet(ctx, "015:fileInfoMap", fileId).Result()

	if fileInfoUnmarshalData != "" {
		var fileInfoData models.RedisFileInfo
		if err := json.Unmarshal([]byte(fileInfoUnmarshalData), &fileInfoData); err != nil {
			return models.RedisFileInfo{}, err
		}
		return fileInfoData, nil
	}
	return models.RedisFileInfo{}, errors.New("db不存在该文件信息")
}

func GetFileId(fileHash string, fileSize int64) string {
	return fmt.Sprintf("%s_%d", fileHash, fileSize)
}

func CreateFileSlice(fileSlice *multipart.FileHeader, fileId string, fileIndex int64) error {
	src, err := fileSlice.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	uploadPath, err := GetUploadDirPath()
	if err != nil {
		return err
	}

	filePath := filepath.Join(uploadPath, fmt.Sprintf("%s_%s", fileId, "_tmp"))
	if err := os.MkdirAll(filePath, 0755); err != nil {
		return err
	}

	dst, err := os.Create(filepath.Join(filePath, fmt.Sprintf("%d", fileIndex)))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

// MergeFileSlices 合并文件切片
func MergeFileSlices(fileId string) error {
	// 获取上传目录路径
	uploadPath, err := GetUploadDirPath()
	if err != nil {
		return err
	}

	// 切片所在目录
	slicePath := filepath.Join(uploadPath, fileId)

	// 最终合并后的文件路径
	finalPath := filepath.Join(uploadPath, fileId+".tmp")

	// 创建最终文件
	destFile, err := os.Create(finalPath)
	if err != nil {
		return fmt.Errorf("创建合并文件失败: %v", err)
	}
	defer destFile.Close()

	// 读取目录下的所有文件
	files, err := os.ReadDir(slicePath)
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
		sliceFiles[index-1] = filepath.Join(slicePath, file.Name())
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

	// 清理切片文件夹
	defer os.RemoveAll(slicePath)

	// 重命名临时文件
	finalFilePath := filepath.Join(uploadPath, fileId)
	defer os.Rename(finalPath, finalFilePath)

	return nil
}
