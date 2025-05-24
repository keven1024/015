package middleware

import (
	"backend/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func LoggerMiddleware() echo.MiddlewareFunc {
	log := utils.GetLogClient()
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info("request",
				zap.String("url", v.URI),
				zap.Int("status", v.Status),
				zap.String("method", v.Method),
				zap.String("ip", v.RemoteIP),
			)
			return nil
		},
	})
}
