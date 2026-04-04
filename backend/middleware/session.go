package middleware

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v5"
)

func SessionMiddleware() echo.MiddlewareFunc {
	store := sessions.NewCookieStore([]byte("secret")) // TODO: 从配置中获取密钥
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("_session_store", store)
			return next(c)
		}
	}
}
