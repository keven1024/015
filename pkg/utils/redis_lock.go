package utils

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislock"
)

const (
	defaultRedisLockPrefix   = "015:lock"
	defaultRedisLockValidity = 15 * time.Second
)

var (
	// ErrRedisLockNotAcquired indicates that TryRedisLock did not obtain the lock.
	ErrRedisLockNotAcquired = errors.New("redis lock not acquired")

	onceLocker       sync.Once
	redisLock        rueidislock.Locker
	newRueidisLocker = rueidislock.NewLocker
)

// RedisLockOption defines the caller-controlled lock settings.
type RedisLockOption struct {
	KeyValidity time.Duration
}

func GetRedisLocker() rueidislock.Locker {
	onceLocker.Do(func() {
		opt, err := rueidis.ParseURL(GetEnv("redis.url"))
		if err != nil {
			panic(err)
		}
		locker, err := newRueidisLocker(rueidislock.LockerOption{
			ClientOption:   opt,
			KeyPrefix:      defaultRedisLockPrefix,
			KeyValidity:    defaultRedisLockValidity,
			KeyMajority:    1,
			NoLoopTracking: false,
			FallbackSETPX:  false,
		})
		if err != nil {
			panic(err)
		}
		redisLock = locker
	})
	return redisLock
}

func baseLocker(ctx context.Context, key string, expired time.Duration) (context.Context, context.CancelFunc, error) {
	if expired <= 0 {
		expired = defaultRedisLockValidity
	}
	locker := GetRedisLocker()
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, expired)
	lockCtx, lockCancel, err := locker.WithContext(timeoutCtx, key)
	if err != nil {
		return nil, nil, err
	}
	return lockCtx, func() {
		lockCancel()
		timeoutCancel()
	}, nil
}

func Locker(key string, expired time.Duration) (context.CancelFunc, error) {
	_, cancel, err := baseLocker(context.Background(), key, expired)
	return cancel, err
}

func WithLocker(ctx context.Context, key string, expired time.Duration, fn func(context.Context) error) error {
	lockCtx, cancel, err := baseLocker(ctx, key, expired)
	if err != nil {
		return err
	}
	defer cancel()
	return fn(lockCtx)
}
