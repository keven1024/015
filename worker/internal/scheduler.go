package internal

import "github.com/hibiken/asynq"

type Schedule struct {
	Cron   string
	Handle *asynq.Task
}

var Schedules = []Schedule{
	{Cron: "0 3 * * *", Handle: asynq.NewTask("file:janitor", nil)},
}
