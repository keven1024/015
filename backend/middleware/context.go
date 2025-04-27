package middleware

import (
	"github.com/labstack/echo/v4"
)

// CustomContext 扩展 echo.Context 以添加自定义功能
type CustomContext struct {
	echo.Context
	Auth interface{}
}

// NewCustomContext 创建自定义上下文的构造函数
func NewCustomContext(c echo.Context) *CustomContext {
	return &CustomContext{
		Context: c,
	}
}

// ContextMiddleware 中间件用于将标准 echo.Context 转换为 CustomContext
func ContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := NewCustomContext(c)
			return next(cc)
		}
	}
}
