package controllers

import (
	"backend/internal/utils"
	statmodel "pkg/models/stat"
	"time"

	"github.com/labstack/echo/v5"
)

const (
	StatDefaultDays = 7
)

type StatRequest struct {
	Type  string `json:"type"`
	Start *int64 `json:"start"`
	End   *int64 `json:"end"`
}

type StatResultEntry struct {
	Timestamp int64 `json:"timestamp"`
	Data      any   `json:"data"`
}

func GetStat(c *echo.Context) error {
	r := new(StatRequest)
	if err := c.Bind(r); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	key := statmodel.StatZSetKey(r.Type)
	if !isStatZSetKey(key) {
		return utils.HTTPErrorHandler(c, ErrInvalidRequest)
	}

	now := time.Now().Unix()
	start := now - int64(StatDefaultDays*24*60*60)
	end := now
	if r.Start != nil {
		start = *r.Start
	}
	if r.End != nil {
		end = *r.End
	}
	if start > end {
		return utils.HTTPErrorHandler(c, ErrInvalidRequest)
	}

	entries, err := statmodel.GetRedisStatRange(key, start, end)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	result := make([]StatResultEntry, 0, len(entries))
	for _, entry := range entries {
		data, err := decodeStatData(key, entry.Data)
		if err != nil {
			return utils.HTTPErrorHandler(c, err)
		}
		result = append(result, StatResultEntry{
			Timestamp: entry.Timestamp,
			Data:      data,
		})
	}

	return utils.HTTPSuccessHandler(c, map[string]any{
		"type":  request.Type,
		"start": startTimestamp,
		"end":   endTimestamp,
		"data":  result,
	})
}

func isStatZSetKey(key statmodel.StatZSetKey) bool {
	switch key {
	case statmodel.FileZSetKey, statmodel.ShareZSetKey, statmodel.DownloadZSetKey, statmodel.WorkerZSetKey:
		return true
	default:
		return false
	}
}

func getStatTimeRange(request *StatRequest) (int64, int64) {
	now := time.Now().Unix()
	startTimestamp := now - int64(StatDefaultDays*24*60*60)
	endTimestamp := now
	if request.Start != nil {
		startTimestamp = *request.Start
	}
	if request.End != nil {
		endTimestamp = *request.End
	}
	return startTimestamp, endTimestamp
}

func decodeStatData(statType statmodel.StatZSetKey, raw string) (any, error) {
	switch statType {
	case statmodel.FileZSetKey:
		return statmodel.JsonFileStatDataToDomain(raw)
	case statmodel.ShareZSetKey:
		return statmodel.JsonShareStatDataToDomain(raw)
	case statmodel.DownloadZSetKey:
		return statmodel.JsonDownloadStatDataToDomain(raw)
	case statmodel.WorkerZSetKey:
		return statmodel.JsonWorkerStatDataToDomain(raw)
	default:
		return nil, ErrInvalidRequest
	}
}
