package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"pkg/utils"

	"github.com/redis/rueidis"
)

type RedisShareInfo struct {
	// Id          string    `json:"id"`
	CreatedAt      int64           `json:"created_at"`
	UpdatedAt      int64           `json:"updated_at"`
	Owner          string          `json:"owner"`
	Type           ShareType       `json:"type"`
	Data           string          `json:"data"` // 分享数据 文件分享为文件id 文本分享为文本内容
	ExpireAt       int64           `json:"expire_time"`
	ViewNum        int64           `json:"download_nums"`
	Password       string          `json:"password"`
	NotifyEmails   []string        `json:"notify_emails"`
	NotifyWebhooks []NotifyWebhook `json:"notify_webhooks"`
	Locale         string          `json:"locale"`
	FileName       string          `json:"file_name"`
	// PickupCode  bool      `json:"pickup_code"`
}

type NotifyWebhook struct {
	URL      string            `json:"url"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers"`
	BodyType string            `json:"bodyType"`
	Body     string            `json:"body"`
}

type ShareType string

const (
	ShareTypeFile ShareType = "file"
	ShareTypeText ShareType = "text"
)

func GetRedisShareInfo(shareId string) (*RedisShareInfo, error) {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	key := fmt.Sprintf("015:shareInfoMap:%s", shareId)
	shareInfoUnmarshalData, err := rdb.Do(ctx, rdb.B().Get().Key(key).Build()).ToString()
	if rueidis.IsRedisNil(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	ttl, _ := rdb.Do(ctx, rdb.B().Ttl().Key(key).Build()).AsInt64()
	var shareInfoData RedisShareInfo

	if err := json.Unmarshal([]byte(shareInfoUnmarshalData), &shareInfoData); err != nil {
		return nil, err
	}
	shareInfoData.ExpireAt = time.Now().Add(time.Duration(ttl) * time.Second).Unix()
	return &shareInfoData, nil
}

func SetRedisShareInfo(shareId string, handler func(shareInfo *RedisShareInfo) *RedisShareInfo) (*RedisShareInfo, error) {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	old_shareInfo, err := GetRedisShareInfo(shareId)
	if err != nil {
		return nil, err
	}
	if old_shareInfo == nil {
		old_shareInfo = &RedisShareInfo{
			CreatedAt: time.Now().Unix(),
		}
	}
	shareInfo := handler(old_shareInfo)
	shareInfo.UpdatedAt = time.Now().Unix()
	jsonData, err := json.Marshal(shareInfo)
	if err != nil {
		return nil, err
	}
	if err := rdb.Do(
		ctx,
		rdb.B().Set().
			Key(fmt.Sprintf("015:shareInfoMap:%s", shareId)).
			Value(string(jsonData)).
			Ex(time.Until(time.Unix(shareInfo.ExpireAt, 0))).
			Build(),
	).Error(); err != nil {
		return nil, err
	}
	return shareInfo, nil
}
