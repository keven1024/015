package models

import (
	"backend/internal/utils"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetRedisPickupData(pickupCode string) (string, error) {
	rdb, ctx := utils.GetRedisClient()
	ShareId, err := rdb.Get(ctx, fmt.Sprintf("015:pickupCode:%s", pickupCode)).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return ShareId, nil
}

func SetRedisPickupData(pickupCode string, shareId string) error {
	rdb, ctx := utils.GetRedisClient()
	_, err := rdb.Set(ctx, fmt.Sprintf("015:pickupCode:%s", pickupCode), shareId, time.Until(time.Now().Add(24*time.Hour))).Result()
	return err
}
