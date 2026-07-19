package share

import (
	"context"
	"fmt"
	"pkg/utils"
	"time"

	"github.com/redis/rueidis"
)

func GetRedisShareInfo(shareId string) (*RedisShareInfo, error) {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	key := fmt.Sprintf("015:shareInfoMap:%s", shareId)
	shareInfoData, err := rdb.Do(ctx, rdb.B().Get().Key(key).Build()).ToString()
	if rueidis.IsRedisNil(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	ttl, _ := rdb.Do(ctx, rdb.B().Ttl().Key(key).Build()).AsInt64()
	shareInfo, err := JsonShareInfoToDomain(shareInfoData)
	if err != nil {
		return nil, err
	}
	shareInfo.ExpireAt = time.Now().Add(time.Duration(ttl) * time.Second).Unix()
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
			Key(fmt.Sprintf("015:shareInfoMap:%s", shareId)).
			Value(jsonData).
			Ex(time.Until(time.Unix(shareInfo.ExpireAt, 0))).
			Build(),
	).Error(); err != nil {
		return nil, err
	}
	return shareInfo, nil
}
