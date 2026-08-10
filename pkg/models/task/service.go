package task

import (
	"context"
	"fmt"
	"pkg/utils"
	"time"

	"github.com/redis/rueidis"
)

func GetRedisTaskInfo(taskId string) (*map[string]any, error) {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	taskInfoData, err := rdb.Do(ctx, rdb.B().Get().Key(fmt.Sprintf("%s:%s", modelName, taskId)).Build()).ToString()
	if rueidis.IsRedisNil(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return JsonTaskInfoToDomain(taskInfoData)
}

func SetRedisTaskInfo(taskId string, taskInfo map[string]any) error {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	jsonData, err := DomainTaskInfoToJson(taskInfo)
	if err != nil {
		return err
	}
	return rdb.Do(
		ctx,
		rdb.B().Set().Key(fmt.Sprintf("%s:%s", modelName, taskId)).Value(jsonData).Ex(time.Hour).Build(),
	).Error()
}
