package models

import (
	"context"
	"encoding/json"

	"pkg/utils"

	"github.com/redis/rueidis"
)

// 统计数据结构
type StatData struct {
	FileSize    int64 `json:"file_size"`    // 文件大小
	FileNum     int64 `json:"file_num"`     // 文件数量
	ShareNum    int64 `json:"share_num"`    // 分享数量
	DownloadNum int64 `json:"download_num"` // 下载数量
}

func GetRedisStat(key string) (*StatData, error) {
	rdb, ctx := utils.GetRedisClient()
	statUnmarshalData, err := rdb.Do(ctx, rdb.B().Hget().Key("015:stat").Field(key).Build()).ToString()
	if rueidis.IsRedisNil(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var stat StatData
	if err := json.Unmarshal([]byte(statUnmarshalData), &stat); err != nil {
		return nil, err
	}
	return &stat, nil
}

func SetRedisStat(key string, handler func(stat *StatData) *StatData) (*StatData, error) {
	var updatedStat *StatData
	err := utils.WithLocker(context.Background(), "015:stat:"+key, 0, func(ctx context.Context) error {
		rdb, _ := utils.GetRedisClient()
		old_stat, err := GetRedisStat(key)
		if err != nil {
			return err
		}
		if old_stat == nil {
			old_stat = &StatData{}
		}
		stat := handler(old_stat)
		jsonData, err := json.Marshal(stat)
		if err != nil {
			return err
		}
		if err := rdb.Do(ctx, rdb.B().Hset().Key("015:stat").FieldValue().FieldValue(key, string(jsonData)).Build()).Error(); err != nil {
			return err
		}
		updatedStat = stat
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updatedStat, nil
}

func GetRedisStatAll() (map[string]string, error) {
	rdb, ctx := utils.GetRedisClient()
	return rdb.Do(ctx, rdb.B().Hgetall().Key("015:stat").Build()).AsStrMap()
}
