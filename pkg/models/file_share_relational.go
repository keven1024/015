package models

import (
	"encoding/json"
	"pkg/utils"

	"github.com/redis/rueidis"
)

func GetRedisFileShareRelational(fileId string) ([]string, error) {
	rdb, ctx := utils.GetRedisClient()
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
	rdb, ctx := utils.GetRedisClient()
	jsonData, _ := json.Marshal(shareIDs)
	return rdb.Do(ctx, rdb.B().Hset().Key("015:fileShareRelational").FieldValue().FieldValue(fileId, string(jsonData)).Build()).Error()
}
