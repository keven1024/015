package models

import (
	"backend/internal/utils"
	"encoding/json"

	"dario.cat/mergo"
	"github.com/redis/go-redis/v9"
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
	statUnmarshalData, err := rdb.HGet(ctx, "015:stat", key).Result()
	if err == redis.Nil {
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

func SetRedisStat(key string, stat StatData) error {
	rdb, ctx := utils.GetRedisClient()
	old_stat, err := GetRedisStat(key)
	if err != nil {
		return err
	}
	if old_stat != nil {
		mergo.Merge(&stat, old_stat)
	}
	jsonData, _ := json.Marshal(stat)
	_, err = rdb.HSet(ctx, "015:stat", key, string(jsonData)).Result()
	return err
}
