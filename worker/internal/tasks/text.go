package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"pkg/models"
	"worker/internal/services"

	"github.com/hibiken/asynq"
)

func TranslateText(ctx context.Context, task *asynq.Task) error {
	var payload TranslateTextTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	result, err := services.TranslateText(payload.Text, payload.Source, payload.Target, payload.Provider)
	if err != nil {
		// 配置缺失或未知提供商属于永久性错误，跳过重试
		if errors.Is(err, services.ErrProviderNotConfigured) || errors.Is(err, services.ErrUnknownProvider) {
			return fmt.Errorf("%w: %w", err, asynq.SkipRetry)
		}
		return err
	}

	return models.SetRedisTaskInfo(task.ResultWriter().TaskID(), map[string]any{
		"status": "success",
		"result": result,
	})
}
