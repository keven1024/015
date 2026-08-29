package share

import (
	"context"
	"fmt"
	"pkg/utils"
	"time"

	"github.com/redis/rueidis"
	"github.com/spf13/cast"
)

func GetRedisShareInfo(shareId string) (*RedisShareInfo, error) {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", modelName, shareId)
	shareInfoData, err := rdb.Do(ctx, rdb.B().Get().Key(key).Build()).ToString()
	if rueidis.IsRedisNil(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	shareInfo, err := JsonShareInfoToDomain(shareInfoData)
	if err != nil {
		return nil, err
	}
	return shareInfo, nil
}

func SetRedisShareInfo(shareId string, handler func(shareInfo *RedisShareInfo) *RedisShareInfo) (*RedisShareInfo, error) {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	oldShareInfo, err := GetRedisShareInfo(shareId)
	if err != nil {
		return nil, err
	}
	if oldShareInfo == nil {
		oldShareInfo = &RedisShareInfo{
			CreatedAt: time.Now().Unix(),
		}
	}
	shareInfo := handler(oldShareInfo)
	shareInfo.UpdatedAt = time.Now().Unix()
	jsonData, err := DomainShareInfoToJson(shareInfo)
	if err != nil {
		return nil, err
	}
	if err := rdb.Do(
		ctx,
		rdb.B().Set().
			Key(fmt.Sprintf("%s:%s", modelName, shareId)).
			Value(jsonData).
			// 分享过期后仍保留下载窗口期，供已签发的 token 下载。
			Ex(time.Until(time.Unix(shareInfo.ExpireAt, 0).Add(
				cast.ToDuration(utils.GetEnvWithDefault("share.download_window", "12")+"h"),
			))).
			Build(),
	).Error(); err != nil {
		return nil, err
	}
	return shareInfo, nil
}

func DeleteRedisShareInfo(shareId string) error {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	return rdb.Do(ctx, rdb.B().Del().Key(fmt.Sprintf("%s:%s", modelName, shareId)).Build()).Error()
}
