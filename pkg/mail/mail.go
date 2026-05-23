package mail

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed out/*.html
var templates embed.FS

func GetMailTemplate(name string) (string, error) {
	data, err := templates.ReadFile(fmt.Sprintf("out/%s.html", name))
	if err != nil {
		return "", fmt.Errorf("mail template %q not found", name)
	}
	return string(data), nil
}

func RenderMailTemplate(name string, vars map[string]string) (string, error) {
	html, err := GetMailTemplate(name)
	if err != nil {
		return "", err
	}
	for k, v := range vars {
		html = strings.ReplaceAll(html, k, v)
	}
	return html, nil
}
