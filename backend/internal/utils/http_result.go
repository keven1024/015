package utils

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type HTTPBaseResponse struct {
	code    int
	message string
	data    map[string]any
}

type HTTPBaseResponseProps func(props *HTTPBaseResponse)

func WithCode(data int) HTTPBaseResponseProps {
	return func(props *HTTPBaseResponse) {
		props.code = data
	}
}

func WithMessage(data string) HTTPBaseResponseProps {
	return func(props *HTTPBaseResponse) {
		props.message = data
	}
}

func WithData(data map[string]any) HTTPBaseResponseProps {
	return func(props *HTTPBaseResponse) {
		props.data = data
	}
}

func HTTPBaseHandler(c *echo.Context, options ...HTTPBaseResponseProps) error {
	props := HTTPBaseResponse{code: http.StatusOK, message: "success", data: map[string]any{}}
	for _, option := range options {
		option(&props)
	}

	return c.JSON(props.code, map[string]any{
		"code":    props.code,
		"message": props.message,
		"data":    props.data,
	})
}

func HTTPSuccessHandler(c *echo.Context, data map[string]any, options ...HTTPBaseResponseProps) error {
	options = append([]HTTPBaseResponseProps{WithData(data)}, options...)
	return HTTPBaseHandler(c, options...)
}

func HTTPErrorHandler(c *echo.Context, err error, options ...HTTPBaseResponseProps) error {
	options = append([]HTTPBaseResponseProps{WithMessage(err.Error()), WithCode(http.StatusBadRequest)}, options...)
	return HTTPBaseHandler(c, options...)
}
