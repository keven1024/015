package models

import (
	"backend/internal/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetRedisTaskInfo(taskId string) (*map[string]any, error) {
	rdb, ctx := utils.GetRedisClient()
	taskInfo := rdb.Get(ctx, fmt.Sprintf("015:taskInfoMap:%s", taskId))
	taskInfoUnmarshalData, err := taskInfo.Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var taskInfoData map[string]any

	if err := json.Unmarshal([]byte(taskInfoUnmarshalData), &taskInfoData); err != nil {
		return nil, err
	}
	return &taskInfoData, nil
}

func SetRedisTaskInfo(taskId string, taskInfo map[string]any) error {
	rdb, ctx := utils.GetRedisClient()
	jsonData, _ := json.Marshal(taskInfo)
	_, err := rdb.Set(ctx, fmt.Sprintf("015:taskInfoMap:%s", taskId), jsonData, time.Hour).Result()
	return err
}
