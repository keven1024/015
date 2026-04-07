package task

import (
	"encoding/json"

	"github.com/labstack/echo/v5"
)

var validProviders = map[string]bool{
	"google":    true,
	"microsoft": true,
	"deeplx":    true,
	"deepseek":  true,
}

var validSources = map[string]bool{
	"auto":  true,
	"zh-CN": true,
	"en":    true,
	"ja":    true,
	"ko":    true,
}

type TranslateTextRequest struct {
	Text     string `json:"text"`
	Source   string `json:"source"`
	Target   string `json:"target"`
	Provider string `json:"provider"`
}

func HandleTextTranslate(c *echo.Context) ([]byte, error) {
	r := new(TranslateTextRequest)
	if err := c.Bind(r); err != nil {
		return nil, err
	}
	if r.Text == "" || r.Target == "" || !validProviders[r.Provider] || !validSources[r.Source] {
		return nil, ErrInvalidRequest
	}
	return json.Marshal(map[string]any{
		"text":     r.Text,
		"source":   r.Source,
		"target":   r.Target,
		"provider": r.Provider,
	})
}
