package models

import (
	"fmt"
	"time"

	"pkg/utils"

	"github.com/redis/rueidis"
)

func GetRedisPickupData(pickupCode string) (string, error) {
	rdb, ctx := utils.GetRedisClient()
	ShareId, err := rdb.Do(ctx, rdb.B().Get().Key(fmt.Sprintf("015:pickupCode:%s", pickupCode)).Build()).ToString()
	if rueidis.IsRedisNil(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return ShareId, nil
}

func SetRedisPickupData(pickupCode string, shareId string) (bool, error) {
	rdb, ctx := utils.GetRedisClient()
	return rdb.Do(
		ctx,
		rdb.B().Set().Key(fmt.Sprintf("015:pickupCode:%s", pickupCode)).Value(shareId).Nx().Ex(24*time.Hour).Build(),
	).AsBool()
}
