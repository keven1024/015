package file

import (
	"context"
	"pkg/utils"
	"time"

	"github.com/redis/rueidis"
	"github.com/spf13/cast"
)

func GetRedisFileInfo(fileId string) (*RedisFileInfo, error) {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	fileInfoData, err := rdb.Do(ctx, rdb.B().Hget().Key(modelName).Field(fileId).Build()).ToString()
	if rueidis.IsRedisNil(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return JsonFileInfoToDomain(fileInfoData)
}

func SetRedisFileInfo(fileId string, handler func(fileInfo *RedisFileInfo) *RedisFileInfo) (*RedisFileInfo, error) {
	var fileInfo *RedisFileInfo
	err := utils.WithLocker(context.Background(), redisLockKey(fileId), 0, func(ctx context.Context) error {
		rdb := utils.GetRedisClient()
		oldFileInfo, err := GetRedisFileInfo(fileId)
		if err != nil {
			return err
		}
		if oldFileInfo == nil {
			oldFileInfo = &RedisFileInfo{
				CreatedAt: time.Now().Unix(),
				Expire:    cast.ToInt64(utils.GetEnvWithDefault("upload.remove_expire", "2")) * 3600,
			}
		}
		fileInfo = handler(oldFileInfo)
		fileInfo.UpdatedAt = time.Now().Unix()
		jsonData, err := DomainFileInfoToJson(fileInfo)
		if err != nil {
			return err
		}
		return rdb.Do(ctx, rdb.B().Hset().Key(modelName).FieldValue().FieldValue(fileId, jsonData).Build()).Error()
	})
	if err != nil {
		return nil, err
	}
	return fileInfo, nil
}

func DelRedisFileInfo(fileId string) error {
	return utils.WithLocker(context.Background(), redisLockKey(fileId), 0, func(ctx context.Context) error {
		rdb := utils.GetRedisClient()
		return rdb.Do(ctx, rdb.B().Hdel().Key(modelName).Field(fileId).Build()).Error()
	})
}

func GetRedisFileInfoAll() (map[string]string, error) {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	return rdb.Do(ctx, rdb.B().Hgetall().Key(modelName).Build()).AsStrMap()
}
