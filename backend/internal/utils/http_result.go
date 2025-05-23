package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Option interface {
	applyTo(*HTTPBaseResponse)
}

type HTTPBaseResponse struct {
	code    int
	message string
	data    map[string]any
}

type HTTPBaseResponseProps func(props *HTTPBaseResponse) error

type WithCode int

func (o WithCode) applyTo(props *HTTPBaseResponse) {
	props.code = int(o)
}

type WithMessage string

func (o WithMessage) applyTo(props *HTTPBaseResponse) {
	props.message = string(o)
}

type WithData map[string]any

func (o WithData) applyTo(props *HTTPBaseResponse) {
	props.data = o
}

func HTTPBaseHandler(c echo.Context, options ...Option) error {
	props := HTTPBaseResponse{code: http.StatusOK, message: "success", data: map[string]any{}}
	for _, option := range options {
		option.applyTo(&props)
	}

	return c.JSON(props.code, map[string]any{
		"code":    props.code,
		"message": props.message,
		"data":    props.data,
	})
}

func HTTPSuccessHandler(c echo.Context, data map[string]any, options ...HTTPBaseResponseProps) error {
	return HTTPBaseHandler(c, WithData(data))
}

func HTTPErrorHandler(c echo.Context, err error, options ...HTTPBaseResponseProps) error {
	return HTTPBaseHandler(c, WithMessage(err.Error()), WithCode(http.StatusBadRequest))
}
