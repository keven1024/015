## Why

当前分享功能支持配置通知（邮件和 Webhook），但通知逻辑从未被实现——下载时不会触发任何通知。用户配置了通知但永远不会收到，导致该功能形同虚设。

## What Changes

- 新增 `NotifyWebhook` 结构体到共享模型层（`pkg/models`），存储 Webhook 的 URL、方法、请求头、Body 类型和 Body 内容
- 扩展 `RedisShareInfo` 模型，增加 `has_notify`、`notify_types`、`notify_emails`、`notify_webhooks` 字段
- 更新后端 `ShareConfig` 请求结构体，与前端实际发送的字段（`notify_types`、`notify_emails`、`notify_webhooks`）对齐
- 在创建分享时将完整通知配置持久化到 Redis
- 在 `VaildateShare`（下载验证）成功后，若分享配置了通知，则入队 `share:notify` 任务
- Worker 实现 `ShareNotify` 处理器：读取分享配置，并发执行所有邮件和 Webhook 通知；仅当全部通知渠道均失败时才返回错误触发重试

## Capabilities

### New Capabilities

- `share-download-notify`: 下载通知能力——分享被下载时，通过邮件和/或 Webhook 向配置的接收方发送通知

### Modified Capabilities

（无既有 Spec 需要变更）

## Impact

- **`pkg/models/share.go`**: 新增 `NotifyWebhook` 类型，扩展 `RedisShareInfo` 字段（影响所有读写 Redis 分享信息的代码）
- **`backend/internal/controllers/share.go`**: `ShareConfig` 字段变更（`notify_email` 重命名为 `notify_emails`，新增 `notify_types`、`notify_webhooks`）
- **`backend/internal/controllers/download.go`**: 新增 asynq 任务入队逻辑，需引入 `encoding/json` 和 `asynq` 依赖
- **`worker/internal/tasks/share.go`**: 实现真实的 `ShareNotify` 函数（当前为空 stub）
- **`worker/internal/tasks/types.go`**: 新增 `ShareNotifyTaskPayload` 结构体
- **`worker/main.go`**: 注册 `share:notify` 路由
- **依赖**: Worker 使用已有的 `go-resty/resty`（Webhook HTTP 请求）和 Go 标准库 `net/smtp`（邮件），无需新增外部依赖
- **配置**: 需新增 SMTP 相关配置项（`smtp.host`、`smtp.port`、`smtp.username`、`smtp.password`、`smtp.from`）
