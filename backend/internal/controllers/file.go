package controllers

import (
	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/utils"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

func CreateUploadTask(c echo.Context) error {
	// cc := c.(*middleware.CustomContext)
	r := new(models.FileInfo)
	if err := c.Bind(r); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	if r.FileSize == 0 || r.MimeType == "" || r.FileHash == "" {
		return utils.HTTPErrorHandler(c, errors.New("上传文件信息不完整"))
	}
	fileId := utils.GetFileId(r.FileHash, r.FileSize)
	fileInfo, _ := models.GetRedisFileInfo(fileId)

	if fileInfo != nil {
		return utils.HTTPSuccessHandler(c, map[string]any{
			"size":       fileInfo.FileSize,
			"mime_type":  fileInfo.MimeType,
			"hash":       fileInfo.FileHash,
			"type":       fileInfo.FileType,
			"expire":     fileInfo.Expire,
			"id":         fileId,
			"chunk_size": fileInfo.ChunkSize,
		})
	}
	ChunkSize := int64(1 * 1024 * 1024)
	if r.FileSize > 500*1024*1024 {
		ChunkSize = r.FileSize / 500
	}
	newFileInfo := models.RedisFileInfo{
		FileType: models.FileTypeInit,
		FileInfo: models.FileInfo{
			FileSize:  r.FileSize,
			MimeType:  r.MimeType,
			FileHash:  r.FileHash,
			ChunkSize: ChunkSize,
		},
		CreatedAt: time.Now().Unix(),
		Expire:    3600,
	}
	err := models.SetRedisFileInfo(fileId, newFileInfo)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	return utils.HTTPSuccessHandler(c, map[string]any{
		"size":       newFileInfo.FileSize,
		"mime_type":  newFileInfo.MimeType,
		"hash":       newFileInfo.FileHash,
		"type":       newFileInfo.FileType,
		"expire":     newFileInfo.Expire,
		"id":         fileId,
		"chunk_size": newFileInfo.ChunkSize,
	})
}

type UploadFileSliceProps struct {
	FileId    string                `form:"id"`
	FileIndex int64                 `form:"index"`
	FileSlice *multipart.FileHeader `form:"file"`
}

func UploadFileSlice(c echo.Context) error {
	r := new(UploadFileSliceProps)
	if err := c.Bind(r); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	if r.FileId == "" || r.FileIndex == 0 || r.FileSlice == nil {
		return utils.HTTPErrorHandler(c, errors.New("上传文件信息不完整"))
	}
	fileInfo, err := models.GetRedisFileInfo(r.FileId)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	now := time.Now().Unix()
	if fileInfo.CreatedAt+fileInfo.Expire < now {
		return utils.HTTPErrorHandler(c, errors.New("上传任务已过期"))
	}

	if fileInfo.FileType != models.FileTypeInit {
		return utils.HTTPErrorHandler(c, errors.New("上传任务状态错误"))
	}
	if r.FileIndex > ((fileInfo.FileSize / fileInfo.ChunkSize) + 1) {
		return utils.HTTPErrorHandler(c, errors.New("文件切片索引错误"))
	}

	if r.FileSlice.Size > fileInfo.ChunkSize {
		return utils.HTTPErrorHandler(c, errors.New("文件切片大小错误"))
	}

	// 打开文件
	file, err := r.FileSlice.Open()
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	defer file.Close()

	if err := services.CreateFileSlice(file, r.FileId, r.FileIndex); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	return utils.HTTPSuccessHandler(c, map[string]any{
		"message": "成功上传",
	})
}

type FinishUploadTaskProps struct {
	FileId string `json:"id"`
}

func FinishUploadTask(c echo.Context) error {
	r := new(FinishUploadTaskProps)
	if err := c.Bind(r); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	if r.FileId == "" {
		return utils.HTTPErrorHandler(c, errors.New("文件ID不能为空"))
	}

	fileInfo, err := models.GetRedisFileInfo(r.FileId)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	if fileInfo.FileType != models.FileTypeInit {
		return utils.HTTPErrorHandler(c, errors.New("上传任务状态错误"))
	}

	now := time.Now().Unix()
	if fileInfo.CreatedAt+fileInfo.Expire < now {
		return utils.HTTPErrorHandler(c, errors.New("上传任务已过期"))
	}

	// 合并文件切片
	uploadPath, _ := services.GetUploadDirPath()
	slicesPath := filepath.Join(uploadPath, fmt.Sprintf("%s_%s", r.FileId, "tmp"))

	// 最终合并后的文件路径
	mergeFilePath := filepath.Join(uploadPath, r.FileId)
	if err := services.MergeFileSlices(slicesPath, mergeFilePath); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	// 计算文件MD5
	file, err := os.Open(mergeFilePath)
	if err != nil {
		file.Close()
		os.Remove(mergeFilePath)
		return utils.HTTPErrorHandler(c, err)
	}

	file_hash, err := utils.GetFileMd5(file)

	if err != nil {
		file.Close()
		os.Remove(mergeFilePath)
		return utils.HTTPErrorHandler(c, err)
	}

	if file_hash != fileInfo.FileHash {
		file.Close()
		os.Remove(mergeFilePath)
		return utils.HTTPErrorHandler(c, errors.New("文件MD5不一致"))
	}
	defer file.Close()
	// 更新文件信息
	models.SetRedisFileInfo(r.FileId, models.RedisFileInfo{
		FileType: models.FileTypeUpload,
	})

	return utils.HTTPSuccessHandler(c, map[string]any{
		"size":      fileInfo.FileSize,
		"mime_type": fileInfo.MimeType,
		"hash":      fileInfo.FileHash,
		"type":      models.FileTypeUpload,
		"id":        r.FileId,
	})
}
