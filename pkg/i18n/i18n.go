package i18n

import (
	"embed"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.toml
var translationFS embed.FS

var (
	bundle             *i18n.Bundle
	defaultLanguage    = language.English
	supportedLanguages = []language.Tag{
		language.English,
		language.Chinese,
		language.French,
	}
)

// 初始化 i18n bundle
func Init() error {
	bundle = i18n.NewBundle(defaultLanguage)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 加载所有翻译文件
	entries, err := translationFS.ReadDir("locales")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".toml") {
			filePath := path.Join("locales", entry.Name())
			data, err := translationFS.ReadFile(filePath)
			if err != nil {
				return err
			}

			_, err = bundle.ParseMessageFileBytes(data, filePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 获取针对特定语言的本地化器
func NewLocalizer(lang string) *i18n.Localizer {
	// 如果传入空字符串，使用默认语言
	if lang == "" {
		return i18n.NewLocalizer(bundle, defaultLanguage.String())
	}

	// 创建包含回退语言的本地化器
	return i18n.NewLocalizer(bundle, lang, defaultLanguage.String())
}

// T 是简化的翻译函数，用于没有参数的简单字符串
func T(lang, messageID string) string {
	localizer := NewLocalizer(lang)
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})

	if err != nil {
		// 记录错误但返回 messageID 作为回退
		return messageID
	}

	return msg
}

// TWithData 用于包含模板数据的翻译
func TWithData(lang, messageID string, templateData map[string]interface{}) string {
	localizer := NewLocalizer(lang)
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})

	if err != nil {
		return messageID
	}

	return msg
}

// TPlural 用于处理复数形式
func TPlural(lang, messageID string, count interface{}, templateData map[string]interface{}) string {
	if templateData == nil {
		templateData = make(map[string]interface{})
	}

	// 确保 templateData 中包含 Count
	templateData["Count"] = count

	localizer := NewLocalizer(lang)
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		PluralCount:  count,
		TemplateData: templateData,
	})

	if err != nil {
		return messageID
	}

	return msg
}
