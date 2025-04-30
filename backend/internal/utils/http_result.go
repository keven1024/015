package utils

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler(c echo.Context, err error, options ...HTTPBaseResponseProps) error {
	return HTTPBaseHandler(c, WithMessage(err.Error()), WithCode(http.StatusBadRequest))
}

type HTTPBaseResponse struct {
	code    *int
	message *string
	data    *map[string]any
}

type HTTPBaseResponseProps func(props *HTTPBaseResponse) error

func WithCode(code int) HTTPBaseResponseProps {
	return func(props *HTTPBaseResponse) error {
		if code < 100 || code > 599 {
			return errors.New("code should be positive")
		}
		props.code = &code
		return nil
	}
}

func WithMessage(message string) HTTPBaseResponseProps {
	return func(props *HTTPBaseResponse) error {
		if message == "" {
			return errors.New("message should not be empty")
		}
		props.message = &message
		return nil
	}
}

func WithData(data map[string]any) HTTPBaseResponseProps {
	return func(props *HTTPBaseResponse) error {
		props.data = &data
		return nil
	}
}
func HTTPBaseHandler(c echo.Context, options ...HTTPBaseResponseProps) error {
	var props HTTPBaseResponse
	for _, option := range options {
		err := option(&props)
		if err != nil {
			return err
		}
	}
	code := http.StatusOK
	if props.code != nil {
		code = *props.code
	}
	message := "success"
	if props.message != nil {
		message = *props.message
	}
	data := map[string]any{}
	if props.data != nil {
		data = *props.data
	}

	return c.JSON(code, map[string]any{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func HTTPSuccessHandler(c echo.Context, data map[string]any, options ...HTTPBaseResponseProps) error {
	return HTTPBaseHandler(c, WithData(data))
}
