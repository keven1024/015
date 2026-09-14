package middleware

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func LoggerMiddleware(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		start := time.Now()
		zap.L().Info("[worker task]::start::", zap.String("type", t.Type()))
		err := h.ProcessTask(ctx, t)
		if err != nil {
			zap.L().Error("[worker task]::error::", zap.String("type", t.Type()), zap.Error(err))
			return err
		}
		zap.L().Info("[worker task]::end::", zap.String("type", t.Type()), zap.Duration("duration", time.Since(start)))
		return nil
	})
}
