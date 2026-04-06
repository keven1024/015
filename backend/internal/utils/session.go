package utils

import (
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v5"
)

func GetSession(c *echo.Context, name string) (*sessions.Session, error) {
	store, err := echo.ContextGet[sessions.Store](c, "_session_store")
	if err != nil {
		return nil, fmt.Errorf("failed to get session store: %w", err)
	}
	return store.Get(c.Request(), name)
}
