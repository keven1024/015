package models

import (
	"backend/internal/utils"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

func GetRedisFileShareRelational(fileId string) ([]string, error) {
	rdb, ctx := utils.GetRedisClient()
	fileShareRelationalUnmarshalData, err := rdb.HGet(ctx, "015:fileShareRelational", fileId).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var shareIDs []string
	if err := json.Unmarshal([]byte(fileShareRelationalUnmarshalData), &shareIDs); err != nil {
		return nil, err
	}
	return shareIDs, nil
}

func SetRedisFileShareRelational(fileId string, shareIDs []string) error {
	rdb, ctx := utils.GetRedisClient()
	jsonData, _ := json.Marshal(shareIDs)
	_, err := rdb.HSet(ctx, "015:fileShareRelational", fileId, string(jsonData)).Result()
	return err
}
