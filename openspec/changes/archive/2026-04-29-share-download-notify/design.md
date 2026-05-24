## Context

项目使用 Go 工作区（go.work）管理多模块：`pkg/models`、`pkg/utils`（共享层）、`backend`（Echo HTTP API）和 `worker`（Asynq 异步任务处理器）。前端已完整实现通知配置 UI（`NotifyConfigField.vue`），调用 share 创建接口时会传入 `notify_types`、`notify_emails`、`notify_webhooks` 等字段。但后端 `ShareConfig` 与前端字段不匹配（`notify_email` 单数 vs `notify_emails` 复数，缺少 `notify_types`、`notify_webhooks`），且这些字段从未被持久化到 Redis，Worker 中的 `ShareNotify` 也只是空 stub。

系统无用户账号体系，"用户" 通过 Session nanoid 标识，没有存储任何语言偏好。前端使用 `@nuxtjs/i18n` 管理 7 种语言（zh-CN、zh-TW、en、ja、ko、fr、de），创建分享时可感知当前语言。

## Goals / Non-Goals

**Goals:**
- 分享被下载时（`VaildateShare` 成功），向配置的邮件地址和 Webhook URL 发送通知
- 通知以 Asynq 异步任务执行，不阻塞下载响应
- 邮件和 Webhook 在同一任务中顺序执行；部分成功即视为任务成功（不重试）；全部失败才返回 error 触发 Asynq 重试
- 邮件内容根据分享创建人的语言（locale）进行本地化，支持 7 种语言

**Non-Goals:**
- 通知发送状态的持久化记录
- 通知的批量/去重处理（每次下载独立触发）
- 支持 TLS 客户端证书校验
- Webhook 内容的 i18n

## Decisions

### D1: 通知配置存储在 `RedisShareInfo` 中

**决策**: 将 `HasNotify`、`NotifyTypes`、`NotifyEmails`、`NotifyWebhooks` 直接存入 `RedisShareInfo`，Worker 读取时通过 `shareId` 查找。

**备选**: 在任务 Payload 中携带完整通知配置。

**理由**: Payload 大小有限制（Asynq 推荐 <1KB），Webhook 配置可能较大。通过 `shareId` 查 Redis 是现有模式（`RemoveShare` 也这么做），保持一致性。

---

### D2: 部分成功即任务成功，全部失败才重试

**决策**: 遍历所有通知目标，收集错误；`成功数 > 0` 则返回 nil；全部失败才返回聚合错误。

**备选**: 任一失败即重试（可能导致重复通知）；全部失败也不重试（丢失通知）。

**理由**: 用户描述"反正通知了，除非全部失败才会重试"——保证至少一个通知到达，避免对已成功的渠道重复发送。

---

### D3: 邮件使用 `github.com/wneessen/go-mail`

**决策**: 使用 `go-mail` 发送邮件，配置通过环境变量注入（`smtp.host`、`smtp.port`、`smtp.username`、`smtp.password`、`smtp.from`）。

**备选**: Go 标准库 `net/smtp`。

**理由**: `go-mail` 同时支持 STARTTLS（587）和 Implicit TLS（465），API 更直观，邮件构造（Subject、Body、Header）更安全，无需手动拼接 raw RFC 2822 格式。解决了标准库不支持 465 端口的问题。

---

### D4: Webhook Body 发送策略

**决策**:
- `none`: 不发送 Body
- `form-data`: 设置 `Content-Type: application/x-www-form-urlencoded`，Body 原样发送
- `raw`: 不设置 Content-Type（由 Headers 中用户自定义），Body 原样发送

**理由**: 前端 Body 字段是 Textarea 原始字符串，不做解析。用户可通过 Headers 自定义 Content-Type 覆盖。

---

### D5: 邮件内容使用 `github.com/nicksnyder/go-i18n/v2` 本地化

**决策**: 在 `worker/internal/i18n/` 目录下维护各语言的 TOML 翻译文件（`active.en.toml`、`active.zh-CN.toml` 等）；Worker 启动时加载所有文件到 `i18n.Bundle`；发送邮件时按 `shareInfo.Locale` 查找对应 `Localizer`，缺失时 fallback 到 `en`。

**备选**: 硬编码英文文本；使用 `text/template` 手动维护多语言字符串。

**理由**: `go-i18n` 是 Go 生态标准 i18n 方案，翻译文件独立于代码，未来新增语言只需添加文件，不改逻辑。

---

### D6: 创建人 locale 由前端随 `CreateShareProps` 传入

**决策**: 在 `ShareConfig`（或 `CreateShareProps` 顶层）新增 `Locale string json:"locale"` 字段，前端在创建分享时将 `useI18n().locale.value` 写入，后端存入 `RedisShareInfo.Locale`。

**备选**:
- 读取 HTTP `Accept-Language` 头：值不可控，浏览器发送的可能与用户在 UI 选择的不同
- 不存 locale，每次查用户偏好：系统无用户账号，无处存储

**理由**: 前端语言是用户显式选择的，最能代表其意图。直接传递是最简单可靠的方案。

## Risks / Trade-offs

- **SMTP 未配置时邮件通知静默跳过** → `smtp.host` 为空时跳过邮件（不计入失败数），记录 warn 日志
- **Webhook 目标不可达** → Asynq 默认重试策略兜底（全部失败时）；单个失败不重试（部分成功语义）
- **Redis 中 share 已过期但任务还在队列** → Worker 查不到 `shareInfo` 时返回 nil（静默跳过），不重试
- **`notify_email`（旧字段名）迁移** → 旧字段不再写入，Redis 历史数据无影响（JSON unmarshal 忽略缺失字段）
- **go-mail 和 go-i18n 为新依赖** → 需更新 `worker/go.mod`；两个库均稳定无争议

## Migration Plan

1. 部署 `pkg/models` 变更（向后兼容，新字段默认零值）
2. 部署 Backend（新字段开始写入 Redis；前端同步传入 `locale`）
3. 部署 Worker（注册 `share:notify` handler，加载 i18n bundle）
4. 无需数据迁移，无需 rollback 特殊操作（feature 是纯增量）

## Open Questions

- 邮件是否需要 HTML 格式，还是纯文本即可？当前设计使用纯文本。
- Webhook 通知是否也需要携带分享元数据（如文件名）在 Body 模板中？当前由用户自定义 Body。
