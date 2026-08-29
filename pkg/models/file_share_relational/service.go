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

func SetRedisFileShareRelational(fileId string, handler func(shareIDs []string) []string) ([]string, error) {
	var shareIDs []string
	err := utils.WithLocker(context.Background(), RedisLockKey(fileId), 0, func(ctx context.Context) error {
		oldShareIDs, err := GetRedisFileShareRelational(fileId)
		if err != nil {
			return err
		}
		shareIDs = handler(oldShareIDs)

		rdb := utils.GetRedisClient()
		if len(shareIDs) == 0 {
			return rdb.Do(ctx, rdb.B().Hdel().Key(modelName).Field(fileId).Build()).Error()
		}
		jsonData, err := DomainFileShareRelationalToJson(shareIDs)
		if err != nil {
			return err
		}
		return rdb.Do(ctx, rdb.B().Hset().Key(modelName).FieldValue().FieldValue(fileId, jsonData).Build()).Error()
	})
	if err != nil {
		return nil, err
	}
	return shareIDs, nil
}
