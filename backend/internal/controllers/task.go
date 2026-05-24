package controllers

import (
	"backend/internal/controllers/task"
	"backend/internal/utils"
	"pkg/models"
	u "pkg/utils"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v5"
)

var handleTaskMap = map[string]func(c *echo.Context) ([]byte, error){
	"image:compress": task.HandleImageCompress,
	"image:convert":  task.HandleImageConvert,
	"text:translate": task.HandleTextTranslate,
}

func CreateTask(c *echo.Context) error {
	taskType := c.Param("type")
	if taskType == "" {
		return utils.HTTPErrorHandler(c, ErrInvalidRequest)
	}
	handleTask, ok := handleTaskMap[taskType]
	if !ok {
		return utils.HTTPErrorHandler(c, ErrTaskNotFound)
	}
	json, err := handleTask(c)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	client := u.GetQueueClient()
	info, err := client.Enqueue(asynq.NewTask(taskType, json), asynq.MaxRetry(3))
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}

	return utils.HTTPSuccessHandler(c, map[string]any{
		"id": info.ID,
	})
}

func GetTask(c *echo.Context) error {
	taskId := c.Param("id")
	if taskId == "" {
		return utils.HTTPErrorHandler(c, ErrInvalidRequest)
	}

	taskInfo, err := models.GetRedisTaskInfo(taskId)
	if err != nil {
		return utils.HTTPErrorHandler(c, err)
	}
	if taskInfo == nil {
		client := u.GetQueueInspector()

		queneTaskInfo, err := client.GetTaskInfo("default", taskId)
		if err != nil {
			return utils.HTTPErrorHandler(c, ErrTaskExpired)
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
