package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

type ImageTaskPayload struct {
	FileId string `json:"file_id"`
}

type CompressImageTaskPayload struct {
	ImageTaskPayload
}

func CompressImage(ctx context.Context, task *asynq.Task) error {
	var payload CompressImageTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	return nil
}
