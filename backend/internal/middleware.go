package internal

import (
	"backend/internal/middleware"

	"github.com/labstack/echo/v5"
)

var Middlewares = []func() echo.MiddlewareFunc{
	middleware.SessionMiddleware,
	middleware.AuthMiddleware,
	middleware.RateLimiterMiddleware,
	middleware.LoggerMiddleware,
}
