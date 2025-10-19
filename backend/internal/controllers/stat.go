package controllers

import (
	"backend/internal/models"
	"backend/internal/utils"
	"encoding/json"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	DateLayout       = "2006-01-02"
	DaysToAnalyze    = 30
	QueueHistoryDays = 30
)

type FileChartData struct {
	FileSize int64  `json:"file_size"`
	FileNum  int64  `json:"file_num"`
	Date     string `json:"date"`
}

type QueueChartData struct {
	Processed int `json:"processed"`
	Failed    int `json:"failed"`
}

func GetStat(c echo.Context) error {
	fileInfoMap, err := models.GetRedisFileInfoAll()
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	now := time.Now()
	cutoffTime := now.Add(-DaysToAnalyze * 24 * time.Hour)

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

		createdAt := time.Unix(fileInfo.CreatedAt, 0)
		if createdAt.After(cutoffTime) {
			dateKey := createdAt.Format(DateLayout)
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

	queueInspector := utils.GetQueueInspector()
	queues, err := queueInspector.History("default", QueueHistoryDays)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	queuesChartData := make(map[string]QueueChartData)
	for _, item := range queues {
		dateKey := item.Date.Format(DateLayout)
		queuesChartData[dateKey] = QueueChartData{
			Processed: item.Processed,
			Failed:    item.Failed,
		}
	}

	return utils.HTTPSuccessHandler(c, map[string]any{
		"chart": map[string]any{
			"storage": fileChartData,
			"queue":   queuesChartData,
		},
	})
}
