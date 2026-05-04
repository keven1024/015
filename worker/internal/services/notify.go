package services

import (
	"fmt"
	"pkg/i18n"
	"pkg/models"
	u "pkg/utils"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
	"github.com/spf13/cast"
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

type EmailTemplateData struct {
	Locale    string
	IP        string
	ShareType models.ShareType
	FileName  string
}

func SendEmail(to string, emailTemplateData EmailTemplateData, options ...mail.Option) error {
	smtp := u.GetEnvMap("smtp")
	if smtp["host"] == "" {
		zap.L().Warn("smtp host is empty, skip share notify email", zap.String("to", to))
		return nil
	}

	host := cast.ToString(smtp["host"])
	if host == "" {
		return fmt.Errorf("smtp.host is required")
	}
	username := cast.ToString(smtp["username"])
	if username == "" {
		return fmt.Errorf("smtp.username is required")
	}
	port := lo.Ternary(cast.ToInt(smtp["port"]) != 0, cast.ToInt(smtp["port"]), mail.DefaultPortSSL)

	templateData := map[string]any{
		"IP":        emailTemplateData.IP,
		"SiteURL":   u.GetEnv("site.url"),
		"ShareType": i18n.T(emailTemplateData.Locale, lo.Ternary(emailTemplateData.ShareType == models.ShareTypeText, "share_type_text", "share_type_file")),
		"FileName":  emailTemplateData.FileName,
	}
	subject := i18n.TWithData(emailTemplateData.Locale, "notify_email_subject", templateData)
	body := i18n.TWithData(emailTemplateData.Locale, "notify_email_body", templateData)
	message := mail.NewMsg()
	if err := message.From(username); err != nil {
		return err
	}
	if err := message.To(to); err != nil {
		return err
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, body)

	options = append([]mail.Option{
		mail.WithPort(port),
		mail.WithUsername(username),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
	}, options...)

	password := cast.ToString(smtp["password"])
	if password != "" {
		options = append(options, mail.WithPassword(password))
	}

	if cast.ToString(smtp["protocol"]) == "ssl" {
		options = append(options, mail.WithSSL())
	}
	if cast.ToString(smtp["protocol"]) == "tls" {
		options = append(options, mail.WithTLSPortPolicy(mail.TLSMandatory))
	}

	client, err := mail.NewClient(host, options...)
	if err != nil {
		return err
	}
	return client.DialAndSend(message)
}
