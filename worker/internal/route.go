package internal

import (
	"context"
	"worker/internal/tasks"

	"github.com/hibiken/asynq"
)

type Route struct {
	Path    string
	Handler func(context.Context, *asynq.Task) error
}

var Routes = []Route{
	{Path: "share:remove", Handler: tasks.RemoveShare},
	{Path: "share:notify", Handler: tasks.ShareNotify},
	{Path: "file:remove", Handler: tasks.RemoveFile},
	{Path: "file:janitor", Handler: tasks.FileJanitor},
	{Path: "image:compress", Handler: tasks.CompressImage},
	{Path: "image:convert", Handler: tasks.ConvertImage},
	// {Path: "text:translate", Handler: tasks.TranslateText},
}
