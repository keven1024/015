package middleware

import (
	"backend/internal/utils"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// CustomMiddleware 创建自定义中间件
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			sess, err := utils.GetSession(c, "session")
			if err != nil {
				return err
			}
			sess.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   86400 * 7,
				HttpOnly: true,
			}
			if sess.Values["auth"] == nil {
				id, err := gonanoid.New()
				if err != nil {
					return err
				}
				sess.Values["auth"] = id
				if err := sess.Save(c.Request(), c.Response()); err != nil {
					return err
				}
			}
			c.Set("auth", sess.Values["auth"])
			return next(c)
		}
	}
}
