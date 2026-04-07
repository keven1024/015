package services

import (
	"testing"
	"worker/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Google Translate ---

func TestTranslateWithGoogle_EnToZhCN(t *testing.T) {
	result, err := services.TranslateText("Hello world", "auto", "zh-CN", "google")
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	t.Logf("en→zh-CN: %s", result)
}

func TestTranslateWithGoogle_ZhCNToEn(t *testing.T) {
	result, err := services.TranslateText("你好世界", "zh-CN", "en", "google")
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	t.Logf("zh-CN→en: %s", result)
}

func TestTranslateWithGoogle_AutoDetect(t *testing.T) {
	result, err := services.TranslateText("こんにちは", "auto", "en", "google")
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	t.Logf("auto(ja)→en: %s", result)
}

func TestTranslateWithGoogle_LongText(t *testing.T) {
	text := "The quick brown fox jumps over the lazy dog. " +
		"This sentence contains every letter of the English alphabet. " +
		"It is commonly used for testing purposes."
	result, err := services.TranslateText(text, "en", "zh-CN", "google")
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	t.Logf("长文本→zh-CN: %s", result)
}

// --- Microsoft Translator ---

func TestTranslateWithMicrosoft_EnToZhCN(t *testing.T) {
	result, err := services.TranslateText("Hello world", "auto", "zh-Hans", "microsoft")
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	t.Logf("en→zh-Hans: %s", result)
}

func TestTranslateWithMicrosoft_ZhCNToEn(t *testing.T) {
	result, err := services.TranslateText("你好世界", "zh-Hans", "en", "microsoft")
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	t.Logf("zh-Hans→en: %s", result)
}

func TestTranslateWithMicrosoft_AutoDetect(t *testing.T) {
	result, err := services.TranslateText("こんにちは", "auto", "en", "microsoft")
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	t.Logf("auto(ja)→en: %s", result)
}

// --- 通用 ---

func TestTranslateText_UnknownProvider(t *testing.T) {
	_, err := services.TranslateText("hello", "en", "zh-CN", "unknown_provider")
	assert.ErrorIs(t, err, services.ErrUnknownProvider)
}
