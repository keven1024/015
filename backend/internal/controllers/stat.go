package controllers

import (
	"backend/internal/utils"
	"pkg/models/stat"
	statmodel "pkg/models/stat"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const (
	StatDefaultDays = 7
)

type StatRequest struct {
	Type  string `json:"type"`
	Start *int64 `json:"start"`
	End   *int64 `json:"end"`
}

type HourlyData struct {
	Count int64  `json:"count"`
	Size  *int64 `json:"size,omitempty"`
}

func GetStat(c *echo.Context) error {
	r := new(StatRequest)
	if err := c.Bind(r); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	key := statmodel.StatZSetKey(r.Type)
	if !lo.Contains([]statmodel.StatZSetKey{
		statmodel.FileZSetKey,
		statmodel.ShareZSetKey,
		statmodel.DownloadZSetKey,
		statmodel.WorkerZSetKey,
	}, key) {
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

	data, err := statmodel.GetRedisStatRange(key, start, end)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	result := lo.MapValues(lo.GroupBy(data, func(item statmodel.RedisStatDataItem) string {
		return cast.ToString(time.Unix(item.Timestamp, 0).Truncate(time.Hour).Unix())
	}), func(values []stat.RedisStatDataItem, _ string) any {

		if key == statmodel.FileZSetKey {
			return map[string]any{
				"count": len(values),
				"size": lo.SumBy(values, func(item statmodel.RedisStatDataItem) int64 {
					data, err := statmodel.JsonFileStatDataToDomain(item.Data)
					if err != nil || data == nil {
						return 0
					}
					return data.FileSize
				}),
			}
		}

		if key == statmodel.WorkerZSetKey {
			return lo.GroupBy(lo.FilterMap(values, func(
				item statmodel.RedisStatDataItem,
				_ int,
			) (*statmodel.WorkerStatData, bool) {
				data, err := statmodel.JsonWorkerStatDataToDomain(item.Data)
				if err != nil || data == nil {
					return nil, false
				}

				return data, true
			}), func(data *statmodel.WorkerStatData) string {
				return data.TaskType
			})
		}

		return map[string]any{
			"count": len(values),
		}
	})

	return utils.HTTPSuccessHandler(c, result)
}
