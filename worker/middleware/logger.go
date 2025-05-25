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
		task_type := zap.String("type", t.Type())
		zap.L().Info("[%q] - 开始处理", task_type)
		err := h.ProcessTask(ctx, t)
		if err != nil {
			zap.L().Error("[%q] - 处理失败 - %v", task_type, zap.Error(err))
			return err
		}
		zap.L().Info("[%q] - 完成处理 | 耗时 %v", task_type, zap.Duration("duration", time.Since(start)))
		return nil
	})
}
