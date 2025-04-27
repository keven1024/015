package middleware

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func SessionMiddleware() echo.MiddlewareFunc {
	return session.Middleware(sessions.NewCookieStore([]byte("secret")))
}
