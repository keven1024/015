package services

import (
	"encoding/json"
	"time"

	"pkg/utils"

	"github.com/hibiken/asynq"
)

func SetFileRemoveTask(fileId string, expire time.Duration) error {
	client := utils.GetQueueClient()
	payload, err := json.Marshal(map[string]any{
		"file_id": fileId,
	})
	if err != nil {
		return err
	}
	_, err = client.Enqueue(
		asynq.NewTask("file:remove", payload),
		asynq.ProcessIn(expire),
		asynq.TaskID("file:remove:"+fileId),
	)
	return err
}
