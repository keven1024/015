package models

import (
	"context"
	"encoding/json"
	"pkg/utils"

	"github.com/redis/rueidis"
)

func GetRedisFileShareRelational(fileId string) ([]string, error) {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	fileShareRelationalUnmarshalData, err := rdb.Do(ctx, rdb.B().Hget().Key("015:fileShareRelational").Field(fileId).Build()).ToString()
	if rueidis.IsRedisNil(err) {
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
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	jsonData, err := json.Marshal(shareIDs)
	if err != nil {
		return err
	}
	return rdb.Do(ctx, rdb.B().Hset().Key("015:fileShareRelational").FieldValue().FieldValue(fileId, string(jsonData)).Build()).Error()
}
