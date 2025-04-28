package controllers

import (
	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/utils"
	"encoding/json"
	"errors"
	"mime/multipart"
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
	rdb, ctx := utils.GetRedisClient()
	fileId := utils.GetFileId(r.FileHash, r.FileSize)
	fileInfo, _ := models.GetRedisFileInfo(fileId)

	if fileInfo != (models.RedisFileInfo{}) {
		return utils.HTTPSuccessHandler(c, map[string]any{
			"size":      fileInfo.FileSize,
			"mime_type": fileInfo.MimeType,
			"hash":      fileInfo.FileHash,
			"type":      fileInfo.FileType,
			"expire":    fileInfo.Expire,
			"id":        fileId,
		})
	}

	newFileInfo := models.RedisFileInfo{
		FileType: models.FileTypeInit,
		FileInfo: models.FileInfo{
			FileSize: r.FileSize,
			MimeType: r.MimeType,
			FileHash: r.FileHash,
		},
		CreatedAt: time.Now().Unix(),
		Expire:    3600,
	}
	jsonData, _ := json.Marshal(newFileInfo)
	rdb.HSet(ctx, "015:fileInfoMap", fileId, string(jsonData)).Result()

	return utils.HTTPSuccessHandler(c, map[string]any{
		"size":      newFileInfo.FileSize,
		"mime_type": newFileInfo.MimeType,
		"hash":      newFileInfo.FileHash,
		"type":      newFileInfo.FileType,
		"expire":    newFileInfo.Expire,
		"id":        fileId,
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
	_, err := models.GetRedisFileInfo(r.FileId)

	if err != nil {
		return utils.HTTPErrorHandler(c, err)
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
	// 合并文件切片
	if err := services.MergeFileSlices(r.FileId); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	// 更新文件状态
	// fileInfo.FileType = models.FileTypeComplete
	if err := services.MergeFileSlices(r.FileId); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	return utils.HTTPSuccessHandler(c, map[string]any{
		"message": "文件上传完成并合并成功",
	})
}
