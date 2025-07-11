package controllers

import (
	"backend/internal/models"
	"backend/internal/utils"
	"backend/middleware"
	"encoding/json"
	"errors"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

type GenCompressImageRequest struct {
	FileId string `json:"file_id"`
}

func GenCompressImage(c echo.Context) error {
	cc := c.(*middleware.CustomContext)
	r := new(GenCompressImageRequest)
	if err := cc.Bind(r); err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	if r.FileId == "" {
		return utils.HTTPErrorHandler(c, errors.New("调用接口参数错误"))
	}
	client := utils.GetQueueClient()
	json, err := json.Marshal(map[string]any{
		"file_id": r.FileId,
	})
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	info, err := client.Enqueue(asynq.NewTask("image:compress", json), asynq.MaxRetry(3))
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	return utils.HTTPSuccessHandler(c, map[string]any{
		"id": info.ID,
	})
}

func GetCompressImage(c echo.Context) error {
	cc := c.(*middleware.CustomContext)
	taskId := cc.Param("id")
	if taskId == "" {
		return utils.HTTPErrorHandler(c, errors.New("调用接口参数错误"))
	}

	taskInfo, err := models.GetRedisTaskInfo(taskId)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	if taskInfo == nil {
		client := utils.GetQueueInspector()

		queneTaskInfo, err := client.GetTaskInfo("default", taskId)
		if err != nil {
			return utils.HTTPErrorHandler(c, errors.New("任务已过期"))
		}
		stateMap := map[asynq.TaskState]string{
			asynq.TaskStateActive:    "processing",
			asynq.TaskStatePending:   "pending",
			asynq.TaskStateScheduled: "scheduled",
			asynq.TaskStateRetry:     "retry",
			asynq.TaskStateArchived:  "archived",
			asynq.TaskStateCompleted: "completed",
		}
		if queneTaskInfo != nil {
			return utils.HTTPSuccessHandler(c, map[string]any{
				"status": stateMap[queneTaskInfo.State],
				"err": map[string]any{
					"message":   queneTaskInfo.LastErr,
					"retry":     queneTaskInfo.Retried,
					"max_retry": queneTaskInfo.MaxRetry,
				},
			})
		}
	}
	return utils.HTTPSuccessHandler(c, *taskInfo)
}
