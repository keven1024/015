package services

import (
	"fmt"
	"pkg/i18n"
	"pkg/models"
	u "pkg/utils"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
	mail "github.com/wneessen/go-mail"
	"go.uber.org/zap"
)

func SendWebhook(webhook models.NotifyWebhook) error {
	method := strings.ToUpper(strings.TrimSpace(webhook.Method))
	if method == "" {
		method = "POST"
	}

	request := resty.New().R()
	if strings.EqualFold(webhook.BodyType, "form-data") {
		request.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	}
	for key, value := range webhook.Headers {
		request.SetHeader(key, value)
	}
	if !strings.EqualFold(webhook.BodyType, "none") && webhook.Body != "" {
		request.SetBody(webhook.Body)
	}

	resp, err := request.Execute(method, webhook.URL)
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 {
		return fmt.Errorf("webhook %s returned status %d", webhook.URL, resp.StatusCode())
	}
	return nil
}

func SendEmail(to string, shareInfo *models.RedisShareInfo, ip string) error {
	host := u.GetEnv("smtp.host")
	if host == "" {
		zap.L().Warn("smtp host is empty, skip share notify email", zap.String("to", to))
		return nil
	}

	username := u.GetEnv("smtp.username")
	password := u.GetEnv("smtp.password")
	from := u.GetEnvWithDefault("smtp.from", username)
	if from == "" {
		return fmt.Errorf("smtp.from or smtp.username is required")
	}

	templateData := map[string]any{
		"IP":        ip,
		"SiteURL":   u.GetEnv("site.url"),
		"ShareType": i18n.T(shareInfo.Locale, lo.Ternary(shareInfo.Type == models.ShareTypeText, "share_type_text", "share_type_file")),
		"FileName":  shareInfo.FileName,
	}
	subject := i18n.TWithData(shareInfo.Locale, "notify_email_subject", templateData)
	body := i18n.TWithData(shareInfo.Locale, "notify_email_body", templateData)
	message := mail.NewMsg()
	if err := message.From(from); err != nil {
		return err
	}
	if err := message.To(to); err != nil {
		return err
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, body)

	port, err := strconv.Atoi(u.GetEnvWithDefault("smtp.port", "587"))
	if err != nil {
		return err
	}

	options := []mail.Option{
		mail.WithPort(port),
	}
	if port == mail.DefaultPortSSL {
		options = append(options, mail.WithSSL())
	} else {
		options = append(options, mail.WithTLSPortPolicy(mail.TLSMandatory))
	}
	if username != "" {
		options = append(options, mail.WithUsername(username), mail.WithPassword(password), mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover))
	}

	client, err := mail.NewClient(host, options...)
	if err != nil {
		return err
	}
	return client.DialAndSend(message)
}
