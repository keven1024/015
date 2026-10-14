package middleware

import (
	"backend/internal/utils"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	sessionMaxAge     = 86400 * 7
	sessionRenewAhead = 4 * 24 * time.Hour
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
				MaxAge:   sessionMaxAge,
				HttpOnly: true,
			}
			now := time.Now()
			expireTime := now.Add(time.Duration(sess.Options.MaxAge) * time.Second).Unix()
			if sess.Values["auth"] == nil {
				id, err := gonanoid.New()
				if err != nil {
					return err
				}
				sess.Values["auth"] = id
				sess.Values["expire_at"] = expireTime
				if err := sess.Save(c.Request(), c.Response()); err != nil {
					return err
				}
			} else {
				expireAt, ok := sess.Values["expire_at"].(int64)
				if !ok || expireAt < now.Add(sessionRenewAhead).Unix() {
					sess.Values["expire_at"] = expireTime
					if err := sess.Save(c.Request(), c.Response()); err != nil {
						return err
					}
				}
			}
			c.Set("auth", sess.Values["auth"])
			return next(c)
		}
	}
}
