## 1. 共享模型层

- [x] 1.1 在 `pkg/models/share.go` 中新增 `NotifyWebhook` 结构体（字段：`URL`、`Method`、`Headers map[string]string`、`BodyType`、`Body`，JSON tag 与前端保持一致：`bodyType`）
- [x] 1.2 移除 `RedisShareInfo.NotifyEmail []string`（旧字段），新增 `HasNotify bool`、`NotifyTypes []string`、`NotifyEmails []string`、`NotifyWebhooks []NotifyWebhook`、`Locale string`

## 2. 后端 —— 创建分享

- [x] 2.1 更新 `backend/internal/controllers/share.go` 中 `ShareConfig` 结构体：移除 `NotifyEmail json:"notify_email"`，新增 `NotifyTypes []string json:"notify_types"`、`NotifyEmails []string json:"notify_emails"`、`NotifyWebhooks []models.NotifyWebhook json:"notify_webhooks"`、`Locale string json:"locale"`
- [x] 2.2 在 `CreateShareInfo` 的 `SetRedisShareInfo` 回调中，将 `HasNotify`、`NotifyTypes`、`NotifyEmails`、`NotifyWebhooks`、`Locale` 从 `r.Config` 赋值到 `shareInfo`

## 3. 后端 —— 下载验证

- [x] 3.1 在 `backend/internal/controllers/download.go` 中，`VaildateShare` 函数锁内的统计更新之后，若 `shareInfo.HasNotify && (len(shareInfo.NotifyEmails) > 0 || len(shareInfo.NotifyWebhooks) > 0)`，则入队 `share:notify` 任务（Payload: `{"share_id": "..."}`)
- [x] 3.2 在 `download.go` 中补充 import：`"encoding/json"` 和 `"github.com/hibiken/asynq"`

## 4. Worker —— 依赖与任务类型

- [x] 4.1 在 `worker/go.mod` 中执行 `go get github.com/wneessen/go-mail` 和 `go get github.com/nicksnyder/go-i18n/v2`，更新 `go.sum`
- [x] 4.2 在 `worker/internal/tasks/types.go` 中新增 `ShareNotifyTaskPayload struct { ShareId string json:"share_id" }`

## 5. Worker —— i18n 翻译文件

- [x] 5.1 创建目录 `worker/internal/i18n/`
- [x] 5.2 创建 `worker/internal/i18n/active.en.toml`：包含邮件 Subject（`notify_email_subject`）和 Body（`notify_email_body`）的英文翻译，Body 中使用模板变量 `{{.ShareType}}`、`{{.FileName}}`
- [x] 5.3 创建另外 6 个语言文件（`active.zh-CN.toml`、`active.zh-TW.toml`、`active.ja.toml`、`active.ko.toml`、`active.fr.toml`、`active.de.toml`），内容与 `en` 保持相同结构，翻译各语言对应文本

## 6. Worker —— 通知处理器

- [x] 6.1 在 `worker/internal/tasks/` 中新建 `notify.go`（或复用 `share.go`），实现 `loadI18nBundle() *i18n.Bundle`：在包级别 `sync.Once` 初始化，加载 `worker/internal/i18n/active.*.toml` 所有文件
- [x] 6.2 实现 `localizeEmail(locale, shareType, fileName string) (subject, body string)`：用 `go-i18n` Localizer 渲染 Subject 和 Body，locale 不存在时 fallback 到 `en`
- [x] 6.3 实现 `sendWebhook(webhook models.NotifyWebhook) error`：使用 resty，按 `BodyType` 设置 Content-Type，附加 Headers，HTTP 状态 >= 400 视为失败
- [x] 6.4 实现 `sendEmail(to string, shareInfo *models.RedisShareInfo) error`：若 `smtp.host` 为空则记录 warn 日志并返回 nil（不计失败）；否则使用 `go-mail` 构造邮件（调用 `localizeEmail` 获取文本）并通过 SMTP 发送
- [x] 6.5 实现 `ShareNotify(ctx, task)` 主函数：解析 Payload → 查 Redis → 遍历 Webhooks + Emails → 收集错误 → 全部失败返回聚合 error，否则返回 nil；share 已过期（`shareInfo == nil`）直接返回 nil

## 7. Worker —— 路由注册

- [x] 7.1 在 `worker/main.go` 的 `mux.HandleFunc` 列表中新增 `mux.HandleFunc("share:notify", tasks.ShareNotify)`
