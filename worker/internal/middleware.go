package internal

import (
	"worker/internal/middleware"

	"github.com/hibiken/asynq"
)

var Middlewares = []func(h asynq.Handler) asynq.Handler{
	middleware.LoggerMiddleware,
}
