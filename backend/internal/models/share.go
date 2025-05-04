package models

import (
	"backend/internal/utils"
	"encoding/json"
	"fmt"
	"time"

	"dario.cat/mergo"
	"github.com/redis/go-redis/v9"
)

type RedisShareInfo struct {
	// Id          string    `json:"id"`
	CreatedAt   int64     `json:"created_at"`
	Owner       string    `json:"owner"`
	Type        ShareType `json:"type"`
	Data        string    `json:"data"` // 分享数据 文件分享为文件id 文本分享为文本内容
	ExpireAt    int64     `json:"expire_time"`
	ViewNum     int64     `json:"download_nums"`
	Password    string    `json:"password"`
	NotifyEmail []string  `json:"notify_email"`
	FileName    string    `json:"file_name"`
	// PickupCode  bool      `json:"pickup_code"`
}

type ShareType string

const (
	ShareTypeFile ShareType = "file"
	ShareTypeText ShareType = "text"
)

func GetRedisShareInfo(shareId string) (*RedisShareInfo, error) {
	rdb, ctx := utils.GetRedisClient()
	shareInfoUnmarshalData, err := rdb.Get(ctx, fmt.Sprintf("015:shareInfoMap:%s", shareId)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var shareInfoData RedisShareInfo
	if err := json.Unmarshal([]byte(shareInfoUnmarshalData), &shareInfoData); err != nil {
		return nil, err
	}
	return &shareInfoData, nil
}

func SetRedisShareInfo(shareId string, shareInfo RedisShareInfo, ttl time.Duration) error {
	rdb, ctx := utils.GetRedisClient()
	old_shareInfo, err := GetRedisShareInfo(shareId)
	if err != nil {
		return err
	}
	if old_shareInfo != nil {
		mergo.Merge(&shareInfo, old_shareInfo)
	}
	jsonData, _ := json.Marshal(shareInfo)
	_, err = rdb.Set(ctx, fmt.Sprintf("015:shareInfoMap:%s", shareId), string(jsonData), ttl).Result()
	return err
}
