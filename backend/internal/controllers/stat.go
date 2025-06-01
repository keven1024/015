package controllers

import (
	"backend/internal/models"
	"backend/internal/utils"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type FileChartData struct {
	FileSize int64  `json:"file_size"`
	FileNum  int64  `json:"file_num"`
	Date     string `json:"date"`
}

func GetStat(c echo.Context) error {
	fileInfoMap, err := models.GetRedisFileInfoAll()
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	fileChartData := make(map[string]FileChartData)
	for _, value := range fileInfoMap {
		var fileInfo models.RedisFileInfo
		err := json.Unmarshal([]byte(value), &fileInfo)
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		if fileInfo.FileType != models.FileTypeUpload {
			continue
		}
		if time.Unix(fileInfo.CreatedAt, 0).After(time.Now().Add(-30 * 24 * time.Hour)) {
			dateKey := time.Unix(fileInfo.CreatedAt, 0).Format("2006-01-02")
			if data, ok := fileChartData[dateKey]; ok {
				fileChartData[dateKey] = FileChartData{
					FileSize: data.FileSize + fileInfo.FileSize,
					FileNum:  data.FileNum + 1,
				}
			} else {
				fileChartData[dateKey] = FileChartData{
					FileSize: fileInfo.FileSize,
					FileNum:  1,
				}
			}
		}
	}
	storageChartData := lo.Times(30, func(i int) FileChartData {
		dateKey := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		if data, ok := fileChartData[dateKey]; ok {
			return FileChartData{
				FileSize: data.FileSize,
				FileNum:  data.FileNum,
				Date:     dateKey,
			}
		}
		return FileChartData{
			FileSize: 0,
			FileNum:  0,
			Date:     dateKey,
		}
	})

	queueInspector := utils.GetQueueInspector()
	queues, err := queueInspector.History("default", 30)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	maxStorageSize, err := utils.GetFileSize(utils.GetEnv("MAX_LOCALSTORAGE_SIZE"))
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	queueData := lo.Map(queues, func(item *asynq.DailyStats, _ int) map[string]any {
		return map[string]any{
			"date":      item.Date.Format("2006-01-02"),
			"processed": item.Processed,
			"failed":    item.Failed,
		}
	})

	return utils.HTTPSuccessHandler(c, map[string]any{
		"version": "0.1.0",
		"max_limit": map[string]any{
			"file_size": maxStorageSize,
		},
		"admin": map[string]any{
			"user_name": utils.GetEnv("ADMIN_NAME"),
			"email":     utils.GetEnv("ADMIN_EMAIL"),
		},
		"chart": map[string]any{
			"storage": storageChartData,
			"queue":   queueData,
		},
	})
}
