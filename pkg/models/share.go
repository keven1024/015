package models

import (
	"encoding/json"
	"fmt"
	"time"

	"pkg/utils"

	"dario.cat/mergo"
	"github.com/redis/rueidis"
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

func SetRedisShareInfo(shareId string, shareInfo RedisShareInfo) error {
	rdb, ctx := utils.GetRedisClient()
	old_shareInfo, err := GetRedisShareInfo(shareId)
	if err != nil {
		return err
	}
	if old_shareInfo != nil {
		mergo.Merge(&shareInfo, old_shareInfo)
	}
	jsonData, _ := json.Marshal(shareInfo)
	return rdb.Do(
		ctx,
		rdb.B().Set().
			Key(fmt.Sprintf("015:shareInfoMap:%s", shareId)).
			Value(string(jsonData)).
			Ex(time.Until(time.Unix(shareInfo.ExpireAt, 0))).
			Build(),
	).Error()
}
