package controllers

import (
	"backend/internal/services"
	"backend/internal/utils"
	"math"
	"mime/multipart"
	"os"
	filemodel "pkg/models/file"
	statmodel "pkg/models/stat"
	s "pkg/services"
	u "pkg/utils"
	"time"

	"github.com/labstack/echo/v5"
)

func CreateUploadTask(c *echo.Context) error {
	// cc := c.(*middleware.CustomContext)
	r := new(filemodel.FileInfo)
	if err := c.Bind(r); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	if r.FileSize == 0 || r.MimeType == "" || r.FileHash == "" {
		return utils.HTTPErrorHandler(c, ErrInvalidRequest)
	}
	fileId := u.GetFileId(r.FileHash, r.FileSize)
	fileInfo, err := filemodel.GetRedisFileInfo(fileId)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	if fileInfo != nil {
		uploadPath, err := u.GetUploadDirPath()
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		sliceList, err := services.GetFileSliceList(fileId, uploadPath)
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		return utils.HTTPSuccessHandler(c, map[string]any{
			"size":       fileInfo.FileSize,
			"mime_type":  fileInfo.MimeType,
			"hash":       fileInfo.FileHash,
			"type":       fileInfo.FileType,
			"expire":     fileInfo.Expire,
			"id":         fileId,
			"chunk_size": fileInfo.ChunkSize,
			"chunks":     sliceList,
		})
	}
	maxStorageSize, err := u.GetFileSize(u.GetEnv("upload.maximum"))
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	fileInfoMap, err := filemodel.GetRedisFileInfoAll()
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	totalSize := int64(0)
	for _, value := range fileInfoMap {
		fileInfo, err := filemodel.JsonFileInfoToDomain(value)
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		totalSize += fileInfo.FileSize
	}
	if totalSize+r.FileSize > int64(maxStorageSize) {
		return utils.HTTPErrorHandler(c, ErrInsufficientStorage)
	}

	ChunkSize := int64(0.25 * 1024 * 1024)
	// 根据文件大小动态调整块大小
	for r.FileSize/ChunkSize > 1000 {
		ChunkSize *= 2
	}
	redisFileInfo, err := filemodel.SetRedisFileInfo(fileId, func(fileInfo *filemodel.RedisFileInfo) *filemodel.RedisFileInfo {
		fileInfo.FileType = filemodel.FileTypeInit
		fileInfo.FileInfo = filemodel.FileInfo{
			FileSize:  r.FileSize,
			MimeType:  r.MimeType,
			FileHash:  r.FileHash,
			ChunkSize: ChunkSize,
		}
		return fileInfo
	})
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	err = s.SetFileRemoveTask(fileId, time.Duration(redisFileInfo.Expire)*time.Second)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	return utils.HTTPSuccessHandler(c, map[string]any{
		"size":       redisFileInfo.FileSize,
		"mime_type":  redisFileInfo.MimeType,
		"hash":       redisFileInfo.FileHash,
		"type":       redisFileInfo.FileType,
		"expire":     redisFileInfo.Expire,
		"id":         fileId,
		"chunk_size": redisFileInfo.ChunkSize,
	})
}

type UploadFileSliceProps struct {
	FileId    string                `form:"id"`
	FileIndex int64                 `form:"index"`
	FileSlice *multipart.FileHeader `form:"file"`
}

func UploadFileSlice(c *echo.Context) error {
	r := new(UploadFileSliceProps)
	if err := c.Bind(r); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	if r.FileId == "" || r.FileIndex == 0 || r.FileSlice == nil {
		return utils.HTTPErrorHandler(c, ErrInvalidRequest)
	}
	fileInfo, err := filemodel.GetRedisFileInfo(r.FileId)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	now := time.Now().Unix()
	if fileInfo.CreatedAt+fileInfo.Expire < now {
		return utils.HTTPErrorHandler(c, ErrUploadTaskExpired)
	}

	if fileInfo.FileType != filemodel.FileTypeInit {
		return utils.HTTPErrorHandler(c, ErrInvalidUploadTaskState)
	}
	if r.FileIndex > ((fileInfo.FileSize / fileInfo.ChunkSize) + 1) {
		return utils.HTTPErrorHandler(c, ErrInvalidFileSliceIndex)
	}

	if r.FileSlice.Size > fileInfo.ChunkSize {
		return utils.HTTPErrorHandler(c, ErrInvalidFileSliceSize)
	}

	// 打开文件
	file, err := r.FileSlice.Open()
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	defer file.Close() //nolint:errcheck

	uploadPath, err := u.GetUploadDirPath()
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	if _, err := services.CreateFileSlice(r.FileId, uploadPath, file, r.FileIndex); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	return utils.HTTPSuccessHandler(c, map[string]any{
		"message": "成功上传",
	})
}

type FinishUploadTaskProps struct {
	FileId string `json:"id"`
}

func FinishUploadTask(c *echo.Context) error {
	r := new(FinishUploadTaskProps)
	if err := c.Bind(r); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	if r.FileId == "" {
		return utils.HTTPErrorHandler(c, ErrInvalidRequest)
	}

	fileInfo, err := filemodel.GetRedisFileInfo(r.FileId)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	if fileInfo.FileType != filemodel.FileTypeInit {
		return utils.HTTPErrorHandler(c, ErrInvalidUploadTaskState)
	}

	now := time.Now().Unix()
	if fileInfo.CreatedAt+fileInfo.Expire < now {
		return utils.HTTPErrorHandler(c, ErrUploadTaskExpired)
	}

	uploadPath, err := u.GetUploadDirPath()
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	fileSliceList, err := services.GetFileSliceList(r.FileId, uploadPath)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	if len(fileSliceList) != int(math.Ceil(float64(fileInfo.FileSize)/float64(fileInfo.ChunkSize))) {
		return utils.HTTPErrorHandler(c, ErrIncompleteFileSlices)
	}

	// 最终合并后的文件路径
	mergeFilePath, err := services.MergeFileSlices(r.FileId, uploadPath)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	// 计算文件SHA1
	file, err := os.Open(mergeFilePath)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	defer file.Close() //nolint:errcheck

	file_hash, err := u.GetFileSHA1(file)
	if err != nil || file_hash != fileInfo.FileHash {
		defer os.Remove(mergeFilePath) //nolint:errcheck
		if err == nil {
			return utils.HTTPErrorHandler(c, ErrFileHashMismatch)
		}
		return utils.HTTPErrorHandler(c, err)
	}

	// 更新文件信息
	fileInfo, err = filemodel.SetRedisFileInfo(r.FileId, func(fileInfo *filemodel.RedisFileInfo) *filemodel.RedisFileInfo {
		fileInfo.FileType = filemodel.FileTypeUpload
		return fileInfo
	})
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	// 统计
	currentDate := time.Now().Format("2006-01-02")
	_, err = statmodel.SetRedisStat(currentDate, func(stat *statmodel.StatData) *statmodel.StatData {
		stat.FileSize += fileInfo.FileSize
		stat.FileNum += 1
		return stat
	})
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	return utils.HTTPSuccessHandler(c, map[string]any{
		"size":      fileInfo.FileSize,
		"mime_type": fileInfo.MimeType,
		"hash":      fileInfo.FileHash,
		"type":      filemodel.FileTypeUpload,
		"id":        r.FileId,
	})
}
