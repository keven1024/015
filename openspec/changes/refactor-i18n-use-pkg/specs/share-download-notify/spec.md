## MODIFIED Requirements

### Requirement: Localized email notification sent on download
`share:notify` Worker 任务 SHALL 对 `notify_emails` 中每个邮件地址通过 SMTP 发送通知邮件，邮件语言使用 `RedisShareInfo.Locale` 对应的翻译，不支持的 locale 回退到英文。支持的语言为：`en`、`de`、`fr`、`ja`、`ko`、`zh-CN`、`zh-TW`。翻译由 `pkg/i18n` 提供，翻译文件以 TOML 格式嵌入二进制。

#### Scenario: Email sent in creator's locale
- **WHEN** SMTP 配置完整（`smtp.host` 非空）且 `shareInfo.Locale="zh-CN"`
- **THEN** Worker 向目标地址发送邮件，邮件 Subject 和 Body 为中文内容

#### Scenario: Email locale fallback to English
- **WHEN** `shareInfo.Locale` 为空或不支持的语言代码
- **THEN** Worker 使用英文模板发送邮件

#### Scenario: SMTP not configured
- **WHEN** `smtp.host` 为空
- **THEN** 邮件通知被跳过（不计入失败数），记录 warn 日志
