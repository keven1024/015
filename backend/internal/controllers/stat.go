package controllers

import (
	"backend/internal/models"
	"backend/internal/utils"
	"encoding/json"

	"github.com/labstack/echo/v4"
)

const (
	DateLayout       = "2006-01-02"
	DaysToAnalyze    = 30
	QueueHistoryDays = 30
)

type StatChartData struct {
	FileSize    int64 `json:"file_size"`
	FileNum     int64 `json:"file_num"`
	ShareNum    int64 `json:"share_num"`
	DownloadNum int64 `json:"download_num"`
}

type QueueChartData struct {
	Processed int `json:"processed"`
	Failed    int `json:"failed"`
}

func GetStat(c echo.Context) error {
	statInfoMap, err := models.GetRedisStatAll()
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	statChartData := make(map[string]StatChartData)
	for key, value := range statInfoMap {
		var statData models.StatData
		err := json.Unmarshal([]byte(value), &statData)
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		statChartData[key] = StatChartData{
			FileSize:    statData.FileSize,
			FileNum:     statData.FileNum,
			ShareNum:    statData.ShareNum,
			DownloadNum: statData.DownloadNum,
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
		if item.Processed == 0 && item.Failed == 0 {
			continue
		}
		queuesChartData[dateKey] = QueueChartData{
			Processed: item.Processed,
			Failed:    item.Failed,
		}
	}

	return utils.HTTPSuccessHandler(c, map[string]any{
		"chart": map[string]any{
			"storage": statChartData,
			"queue":   queuesChartData,
		},
	})
}
