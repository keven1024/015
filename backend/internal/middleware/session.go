package middleware

import (
	u "pkg/utils"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v5"
)

func SessionMiddleware() echo.MiddlewareFunc {
	store := sessions.NewCookieStore([]byte(u.GetConfig().Share.SessionSalt))
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("_session_store", store)
			return next(c)
		}
	}
}
