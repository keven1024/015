package middleware

import (
	"context"
	"time"

	statmodel "pkg/models/stat"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func LoggerMiddleware(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		start := time.Now()
		taskID, _ := asynq.GetTaskID(ctx)
		taskType := t.Type()
		zap.L().Info("[worker task]::start::", zap.String("id", taskID), zap.String("type", taskType))
		err := h.ProcessTask(ctx, t)
		status := "processed"
		if err != nil {
			status = "failed"
			zap.L().Error("[worker task]::error::", zap.String("id", taskID), zap.String("type", taskType), zap.Error(err))
		} else {
			zap.L().Info("[worker task]::end::", zap.String("id", taskID), zap.String("type", taskType), zap.Duration("duration", time.Since(start)))
		}
		if statErr := appendWorkerStat(taskID, taskType, status); statErr != nil {
			zap.L().Warn("[worker task]::stat::failed::", zap.String("id", taskID), zap.String("type", taskType), zap.Error(statErr))
		}
		if err != nil {
			return err
		}
		return nil
	})
}

func appendWorkerStat(taskID string, taskType string, status string) error {
	statData, err := statmodel.DomainWorkerStatDataToJson(&statmodel.WorkerStatData{
		TaskID:   taskID,
		TaskType: taskType,
		Status:   status,
	})
	if err != nil {
		return err
	}
	return statmodel.AppendRedisStat(statmodel.WorkerZSetKey, time.Now().Unix(), statData)
}
