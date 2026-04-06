package models

import (
	"encoding/json"
	"fmt"
	"time"

	"pkg/utils"

	"github.com/redis/rueidis"
)

func GetRedisTaskInfo(taskId string) (*map[string]any, error) {
	rdb, ctx := utils.GetRedisClient()
	taskInfoUnmarshalData, err := rdb.Do(ctx, rdb.B().Get().Key(fmt.Sprintf("015:taskInfoMap:%s", taskId)).Build()).ToString()
	if rueidis.IsRedisNil(err) {
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
	jsonData, err := json.Marshal(taskInfo)
	if err != nil {
		return err
	}
	return rdb.Do(
		ctx,
		rdb.B().Set().Key(fmt.Sprintf("015:taskInfoMap:%s", taskId)).Value(string(jsonData)).Ex(time.Hour).Build(),
	).Error()
}
