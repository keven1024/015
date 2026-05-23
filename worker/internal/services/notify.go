package services

import (
	"fmt"
	"net/url"
	"pkg/i18n"
	pkgmail "pkg/mail"
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
	Region    string
	FileName  string
	ShareType models.ShareType
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
	port := mail.DefaultPortSSL
	if smtpPort := cast.ToInt(smtp["port"]); smtpPort != 0 {
		port = smtpPort
	}

	p, err := url.Parse(u.GetEnv("site.url"))
	subject := i18n.TWithData(emailTemplateData.Locale, "notify_email_subject", map[string]any{
		"SiteURL": p.Host,
	})
	htmlBody, err := pkgmail.RenderMailTemplate("pull-notify", map[string]string{
		"EMAIL-TITLE":        subject,
		"EMAIL-INTRO":        i18n.T(emailTemplateData.Locale, "notify_email_intro"),
		"EMAIL-FILEICON":     lo.Ternary(emailTemplateData.ShareType == models.ShareTypeText, "spiral_notepad", "file_folder"),
		"EMAIL-FILENAME":     emailTemplateData.FileName,
		"EMAIL-IP":           emailTemplateData.IP,
		"EMAIL-REGION":       emailTemplateData.Region,
		"EMAIL-LABEL-REGION": i18n.T(emailTemplateData.Locale, "notify_email_label_region"),
	})
	if err != nil {
		return err
	}
	message := mail.NewMsg()
	if err := message.From(username); err != nil {
		return err
	}
	if err := message.To(to); err != nil {
		return err
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextHTML, htmlBody)

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
