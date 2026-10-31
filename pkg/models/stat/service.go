package stat

import (
	"context"
	"pkg/utils"
	"strconv"

	"github.com/redis/rueidis"
	"github.com/samber/lo"
)

func getRedisStatKey(key StatZSetKey) string {
	return modelName + ":" + string(key)
}

// AppendRedisStat 将一条统计数据追加到指定的 zset 中。
// timestamp 作为 zset score，data 作为 zset member，通常传入对应结构体的 JSON。
func AppendRedisStat(key StatZSetKey, timestamp int64, data string) error {
	rdb := utils.GetRedisClient()
	return rdb.Do(
		context.Background(),
		rdb.B().Zadd().
			Key(getRedisStatKey(key)).
			ScoreMember().
			ScoreMember(float64(timestamp), data).
			Build(),
	).Error()
}

// GetRedisStatRange 获取指定时间戳范围内的统计数据及其 zset score。
func GetRedisStatRange(key StatZSetKey, startTimestamp int64, endTimestamp int64) ([]RedisStatDataItem, error) {
	rdb := utils.GetRedisClient()
	values, err := rdb.Do(
		context.Background(),
		rdb.B().Zrangebyscore().
			Key(getRedisStatKey(key)).
			Min(strconv.FormatInt(startTimestamp, 10)).
			Max(strconv.FormatInt(endTimestamp, 10)).
			Withscores().
			Build(),
	).AsZScores()
	if err != nil {
		return nil, err
	}
	entries := lo.Map(values, func(value rueidis.ZScore, _ int) RedisStatDataItem {
		return RedisStatDataItem{
			Timestamp: int64(value.Score),
			Data:      value.Member,
		}
	})
	return entries, nil
}
