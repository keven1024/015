## 为什么

Webhook 通知目前无法配置请求体，导致所有需要 Payload 的 POST/PUT/PATCH Webhook 均无法使用。此外，Webhook 配置界面的 curl 导入功能已部分实现但不可用——解析结果仅输出到控制台，从未写入表单。

## 变更内容

- 在前端 `WebhookItem` 类型中新增 `body` 字段，在 Webhook 高级设置面板中以 textarea 形式展示。
- 激活 curl 导入处理逻辑：用户在 URL 字段粘贴 `curl` 命令后，自动将 `url`、`method`、`headers`、`body` 填入表单，并展开高级面板。
- 在所有现有 locale 文件中新增 `webhookBody` i18n 翻译键。

## 功能模块

### 新增功能
- `webhook-body-field`：Webhook 通知的请求体配置——用户可为每条 Webhook 指定随请求发送的原始 body 字符串。

### 修改功能
- `share-download-notify`：Webhook 需求扩展——`notify_webhooks` 条目现在携带可选的 `body` 字段；Worker 在该字段有值时 SHALL 将其作为 HTTP 请求体发送。

## 影响范围

- **前端**：`front/components/Preprocessing/NotifyConfigField.vue` — 接口定义、模板、curl 处理函数、默认值。
- **前端 i18n**：`front/` 下所有包含 `webhookHeaders` 键的 locale 文件。
- **后端模型**：`pkg/models/share.go` — `NotifyWebhook` 结构体已有 `Body` 字段，无需变更。
- **Worker**：`worker/internal/tasks/notify.go` — `sendWebhook` 已支持 `Body`，无需变更。
