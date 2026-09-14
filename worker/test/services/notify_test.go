package services

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"pkg/i18n"
	"pkg/models"
	sharemodel "pkg/models/share"
	"pkg/utils"
	"worker/internal/services"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wneessen/go-mail"
)

func TestMain(m *testing.M) {
	utils.SetConfig(models.Config{})
	if err := i18n.Init(); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

// ============= SendWebhook =============

func TestSendWebhook_DefaultMethodIsPost(t *testing.T) {
	var gotMethod string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	err := services.SendWebhook(sharemodel.NotifyWebhook{URL: ts.URL})
	require.NoError(t, err)
	assert.Equal(t, "POST", gotMethod)
}

func TestSendWebhook_CustomMethod(t *testing.T) {
	var gotMethod string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	err := services.SendWebhook(sharemodel.NotifyWebhook{URL: ts.URL, Method: "  put  "})
	require.NoError(t, err)
	assert.Equal(t, "PUT", gotMethod)
}

func TestSendWebhook_CustomHeaders(t *testing.T) {
	var gotHeader string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Custom-Token")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	err := services.SendWebhook(sharemodel.NotifyWebhook{
		URL:     ts.URL,
		Headers: map[string]string{"X-Custom-Token": "secret123"},
	})
	require.NoError(t, err)
	assert.Equal(t, "secret123", gotHeader)
}

func TestSendWebhook_FormDataSetsContentType(t *testing.T) {
	var gotContentType string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	err := services.SendWebhook(sharemodel.NotifyWebhook{
		URL:      ts.URL,
		BodyType: "form-data",
		Body:     "key=value",
	})
	require.NoError(t, err)
	assert.Equal(t, "application/x-www-form-urlencoded", gotContentType)
}

func TestSendWebhook_BodyNoneSkipsBody(t *testing.T) {
	var gotBodyLen int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBodyLen = len(b)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	err := services.SendWebhook(sharemodel.NotifyWebhook{
		URL:      ts.URL,
		BodyType: "none",
		Body:     "should-not-be-sent",
	})
	require.NoError(t, err)
	assert.Zero(t, gotBodyLen)
}

func TestSendWebhook_4xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	err := services.SendWebhook(sharemodel.NotifyWebhook{URL: ts.URL})
	assert.Error(t, err)
}

func TestSendWebhook_5xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	err := services.SendWebhook(sharemodel.NotifyWebhook{URL: ts.URL})
	assert.Error(t, err)
}

func TestSendWebhook_InvalidURLReturnsError(t *testing.T) {
	err := services.SendWebhook(sharemodel.NotifyWebhook{URL: "http://127.0.0.1:1"})
	assert.Error(t, err)
}

// ============= SendEmail =============

func TestSendEmail_HappyPath(t *testing.T) {
	t.Cleanup(func() { utils.SetConfig(models.Config{}) })
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       true,
		LogServerActivity: true,
	})
	require.NoError(t, server.Start())
	t.Cleanup(func() { _ = server.Stop() })

	utils.SetConfig(models.Config{SMTP: &models.SMTPConfig{
		Host:     "127.0.0.1",
		Port:     uint(server.PortNumber()),
		Username: "sender@example.com",
	}})

	err := services.SendEmail("recipient@example.com", services.EmailTemplateData{
		Locale:    "en",
		ShareType: sharemodel.ShareTypeText,
		FileName:  "report.pdf",
		IP:        "1.2.3.4",
	}, mail.WithHELO("localhost"), mail.WithTLSPolicy(mail.NoTLS), mail.WithSMTPAuth(mail.SMTPAuthNoAuth))
	require.NoError(t, err)
}
