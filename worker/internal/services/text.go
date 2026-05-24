package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pkg/utils"
	"strings"
	"worker/internal/services/llm"

	"github.com/go-resty/resty/v2"
)

var (
	ErrProviderNotConfigured = errors.New("ProviderNotConfigured")
	ErrUnknownProvider       = errors.New("UnknownProvider")
)

const microsoftAuthEndpoint = "https://edge.microsoft.com/translate/auth"

func TranslateText(text, from, to, provider string) (string, error) {
	switch provider {
	case "google":
		return translateWithGoogle(text, from, to)
	case "microsoft":
		return translateWithMicrosoft(text, from, to)
	case "deeplx":
		return translateWithDeepLX(text, from, to)
	case "deepseek":
		return translateWithLLM(text, from, to)
	default:
		return "", fmt.Errorf("%w: %s", ErrUnknownProvider, provider)
	}
}

// translateWithGoogle 调用 Google Translate 非官方 API。
// 响应格式：[[["翻译结果","原文",...]], null, "en"]
func translateWithGoogle(text, from, to string) (string, error) {
	endpoint := fmt.Sprintf(
		"https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&q=%s",
		url.QueryEscape(from), url.QueryEscape(to), url.QueryEscape(text),
	)
	resp, err := http.Get(endpoint) //nolint:noctx
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var raw []any
	if err := json.Unmarshal(body, &raw); err != nil {
		return "", err
	}
	if len(raw) == 0 {
		return "", errors.New("GoogleTranslateUnexpectedResponse")
	}
	segments, ok := raw[0].([]any)
	if !ok {
		return "", errors.New("GoogleTranslateUnexpectedResponse")
	}
	var parts []string
	for _, seg := range segments {
		pair, ok := seg.([]any)
		if !ok || len(pair) == 0 {
			continue
		}
		if s, ok := pair[0].(string); ok && s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, ""), nil
}

// translateWithMicrosoft 调用微软 Edge Translator 免费端点（无需 API Key）。
// 响应格式：[{"translations":[{"text":"翻译结果","to":"zh-Hans"}]}]
func translateWithMicrosoft(text, from, to string) (string, error) {
	params := url.Values{}
	params.Set("api-version", "3.0")
	params.Set("to", to)
	if from != "auto" {
		params.Set("from", from)
	}
	endpoint := "https://api-edge.cognitive.microsofttranslator.com/translate?" + params.Encode()

	reqBody, err := json.Marshal([]map[string]string{{"Text": text}})
	if err != nil {
		return "", err
	}
	token, err := getMicrosoftToken()
	if err != nil {
		return "", err
	}

	client := resty.New()
	req := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody)

	req.SetAuthToken(token)

	resp, err := req.Post(endpoint)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("MicrosoftTranslateUnexpectedResponse")
	}
	if resp.IsError() {
		return "", fmt.Errorf("MicrosoftTranslateRequestFailed: status=%d body=%s", resp.StatusCode(), resp.String())
	}

	body := resp.Body()
	if len(body) == 0 {
		return "", errors.New("MicrosoftTranslateUnexpectedResponse")
	}

	var result []struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if len(result) == 0 || len(result[0].Translations) == 0 {
		return "", errors.New("MicrosoftTranslateUnexpectedResponse")
	}
	return result[0].Translations[0].Text, nil
}

func getMicrosoftToken() (string, error) {
	resp, err := resty.New().R().Get(microsoftAuthEndpoint)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("MicrosoftAuthUnexpectedResponse")
	}
	if resp.IsError() {
		return "", fmt.Errorf("MicrosoftAuthRequestFailed: status=%d body=%s", resp.StatusCode(), resp.String())
	}

	token := strings.TrimSpace(resp.String())
	if token == "" {
		return "", errors.New("MicrosoftAuthUnexpectedResponse")
	}
	return token, nil
}

type deeplxRequest struct {
	Text       string `json:"text"`
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
}

type deeplxResponse struct {
	Data string `json:"data"`
}

// translateWithDeepLX 调用自托管 DeepLX 服务，端点通过 deeplx.endpoint 配置。
func translateWithDeepLX(text, from, to string) (string, error) {
	endpoint := utils.GetEnv("deeplx.endpoint")
	if endpoint == "" {
		return "", ErrProviderNotConfigured
	}
	reqBody, err := json.Marshal(deeplxRequest{Text: text, SourceLang: from, TargetLang: to})
	if err != nil {
		return "", err
	}
	resp, err := http.Post(endpoint, "application/json", strings.NewReader(string(reqBody))) //nolint:noctx
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result deeplxResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// translateWithLLM 通过 OpenAI 兼容 API（如 DeepSeek）进行 AI 翻译。
func translateWithLLM(text, from, to string) (string, error) {
	endpoint := utils.GetEnv("llm.endpoint")
	apiKey := utils.GetEnv("llm.api_key")
	_, err := llm.NewClient(endpoint, apiKey)
	if err != nil {
		return "", ErrProviderNotConfigured
	}
	model := utils.GetEnvWithDefault("llm.model", "deepseek-chat")
	system := fmt.Sprintf(
		"You are a professional translator. Translate the given text from %s to %s. "+
			"Output only the translated text, with no explanations or surrounding quotes.",
		from, to,
	)
	return llm.Chat(context.Background(), system, text, model)
}
