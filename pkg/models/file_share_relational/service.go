package file_share_relational

import (
	"context"
	"pkg/utils"

	"github.com/redis/rueidis"
)

func GetRedisFileShareRelational(fileId string) ([]string, error) {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	shareIDsData, err := rdb.Do(ctx, rdb.B().Hget().Key(modelName).Field(fileId).Build()).ToString()
	if rueidis.IsRedisNil(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return JsonFileShareRelationalToDomain(shareIDsData)
}

func SetRedisFileShareRelational(fileId string, shareIDs []string) error {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	jsonData, err := DomainFileShareRelationalToJson(shareIDs)
	if err != nil {
		return err
	}
	return rdb.Do(ctx, rdb.B().Hset().Key(modelName).FieldValue().FieldValue(fileId, jsonData).Build()).Error()
}
