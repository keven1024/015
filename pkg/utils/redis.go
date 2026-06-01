package utils

import (
	"fmt"
	"log/slog"
	"math"
	"time"
	"github.com/redis/rueidis"
)

var rdb rueidis.Client

const (
	redisMaxRetries    = 10
	redisBaseDelay     = 300 * time.Millisecond
	redisBackoffFactor = 2.7
	redisMaxDelay      = 15 * time.Second
)

func InitRedis() error {
	opt, err := rueidis.ParseURL(GetEnv("redis.url"))
	if err != nil {
		return fmt.Errorf("invalid redis url: %w", err)
	}

	var lastErr error
	for attempt := range redisMaxRetries {
		var client rueidis.Client
		client, lastErr = rueidis.NewClient(opt)
		if lastErr == nil {
			rdb = client
			return nil
		}
		if attempt == redisMaxRetries-1 {
			break
		}
		delay := time.Duration(math.Min(
			float64(redisBaseDelay)*math.Pow(redisBackoffFactor, float64(attempt)),
			float64(redisMaxDelay),
		))
		slog.Warn("redis connection failed, retrying",
			"attempt", attempt+1,
			"maxRetries", redisMaxRetries,
			"retryIn", delay.String(),
			"error", lastErr,
		)
		time.Sleep(delay)
	}
	return fmt.Errorf("redis connection failed after %d attempts: %w", redisMaxRetries, lastErr)
}

func GetRedisClient() rueidis.Client {
	return rdb
}
