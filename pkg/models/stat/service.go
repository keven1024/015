package stat

import (
	"context"
	"pkg/utils"

	"github.com/redis/rueidis"
)

func GetRedisStat(key string) (*StatData, error) {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	statData, err := rdb.Do(ctx, rdb.B().Hget().Key(modelName).Field(key).Build()).ToString()
	if rueidis.IsRedisNil(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return JsonStatDataToDomain(statData)
}

func SetRedisStat(key string, handler func(stat *StatData) *StatData) (*StatData, error) {
	var updatedStat *StatData
	err := utils.WithLocker(context.Background(), modelName+":"+key, 0, func(ctx context.Context) error {
		rdb := utils.GetRedisClient()
		oldStat, err := GetRedisStat(key)
		if err != nil {
			return err
		}
		if oldStat == nil {
			oldStat = &StatData{}
		}
		stat := handler(oldStat)
		jsonData, err := DomainStatDataToJson(stat)
		if err != nil {
			return err
		}
		if err := rdb.Do(ctx, rdb.B().Hset().Key(modelName).FieldValue().FieldValue(key, jsonData).Build()).Error(); err != nil {
			return err
		}
		updatedStat = stat
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updatedStat, nil
}

func GetRedisStatAll() (map[string]string, error) {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	return rdb.Do(ctx, rdb.B().Hgetall().Key(modelName).Build()).AsStrMap()
}
