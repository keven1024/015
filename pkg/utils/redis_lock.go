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

func baseLocker(ctx context.Context, key string) (context.Context, context.CancelFunc, error) {
	locker := GetRedisLocker()
	return locker.WithContext(ctx, key)
}

func Locker(key string) (context.CancelFunc, error) {
	_, cancel, err := baseLocker(context.Background(), key)
	if err != nil {
		return nil, err
	}
	return cancel, nil
}

func WithLocker(ctx context.Context, key string, fn func(context.Context) error) error {
	lockCtx, cancel, err := baseLocker(ctx, key)
	if err != nil {
		return err
	}
	defer cancel()
	return fn(lockCtx)
}

func TryWithLocker(ctx context.Context, key string, fn func(context.Context) error) error {
	locker := GetRedisLocker()
	lockCtx, cancel, err := locker.TryWithContext(ctx, key)
	if err != nil {
		if errors.Is(err, rueidislock.ErrNotLocked) {
			return ErrRedisLockNotAcquired
		}
		return err
	}
	defer cancel()
	return fn(lockCtx)
}
