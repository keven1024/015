package middleware

import (
	"net/http"
	"slices"
	"time"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

type RateSkiper struct {
	Path   string
	Method string
}

// 并发上传和下载接口绕过速率限制
var RateSkipList = []RateSkiper{
	{Path: "/file/slice", Method: "POST"},
	{Path: "/download", Method: "GET"},
}

func RateLimiterMiddleware() echo.MiddlewareFunc {
	config := echo_middleware.RateLimiterConfig{
		Skipper: func(e echo.Context) bool {
			path := e.Path()
			r := e.Request()
			return slices.Contains(RateSkipList, RateSkiper{Path: path, Method: r.Method})
		},
		Store: echo_middleware.NewRateLimiterMemoryStoreWithConfig(
			echo_middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(10), Burst: 30, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}
	return echo_middleware.RateLimiterWithConfig(config)
}
