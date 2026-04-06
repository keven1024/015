package services

import (
	"encoding/json"
	"time"

	"pkg/utils"

	"github.com/hibiken/asynq"
)

func SetFileRemoveTask(fileId string, expire time.Duration) error {
	client := utils.GetQueueClient()
	json, err := json.Marshal(map[string]any{
		"file_id": fileId,
	})
	if err != nil {
		return err
	}
	_, err = client.Enqueue(asynq.NewTask("file:remove", json), asynq.ProcessIn(expire))
	return err
}
