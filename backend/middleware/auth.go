package middleware

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// CustomMiddleware 创建自定义中间件
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get("session", c)
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

			cc := c.(*CustomContext)
			cc.Auth = sess.Values["auth"]
			// 将自定义上下文传递给下一个处理器
			return next(cc)
		}
	}
}
