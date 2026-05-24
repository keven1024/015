## ADDED Requirements

### Requirement: Notify config and creator locale persisted on share creation
创建分享时，系统 SHALL 将完整通知配置（`has_notify`、`notify_types`、`notify_emails`、`notify_webhooks`）以及创建人当前语言（`locale`）持久化到 Redis `RedisShareInfo` 中。

#### Scenario: Share created with email and webhook notify config and locale
- **WHEN** 客户端 POST `/share` 并携带 `config.has_notify=true`、`config.notify_types=["email","webhook"]`、`config.notify_emails`、`config.notify_webhooks`、`config.locale="zh-CN"`
- **THEN** Redis 中存储的 `RedisShareInfo` 包含完整的通知配置字段及 `locale="zh-CN"`

#### Scenario: Share created without notify
- **WHEN** 客户端 POST `/share` 并携带 `config.has_notify=false`
- **THEN** Redis 中 `has_notify=false`，通知相关字段为空/零值

---

### Requirement: Notify task enqueued on download validation
当分享下载验证（`VaildateShare`）成功且分享配置了通知时，系统 SHALL 向 Asynq 队列入队一个 `share:notify` 任务。

#### Scenario: Download validated with notify enabled
- **WHEN** 用户成功通过 `VaildateShare` 验证（密码正确、下载次数充足）且 `shareInfo.HasNotify=true` 且存在至少一个通知目标
- **THEN** 系统入队一条 `share:notify` 任务，Payload 包含 `share_id`
- **THEN** 下载 token 正常返回，响应不受通知入队影响

#### Scenario: Download validated with notify disabled
- **WHEN** `shareInfo.HasNotify=false`
- **THEN** 不入队任何通知任务

---

### Requirement: Webhook notification sent on download
`share:notify` Worker 任务 SHALL 对 `notify_webhooks` 中每个配置项发起 HTTP 请求。

#### Scenario: POST webhook with form-data body
- **WHEN** Webhook `method=POST`、`bodyType=form-data`、`body="k=v"`
- **THEN** Worker 向 `url` 发起 POST 请求，`Content-Type: application/x-www-form-urlencoded`，Body 为 `k=v`

#### Scenario: POST webhook with raw body
- **WHEN** Webhook `method=POST`、`bodyType=raw`
- **THEN** Worker 向 `url` 发起 POST 请求，Body 原样发送，不强制设置 Content-Type

#### Scenario: Webhook with custom headers
- **WHEN** Webhook `headers` 中包含自定义请求头
- **THEN** Worker 将这些 Header 附加到请求中

#### Scenario: Webhook returns 4xx/5xx
- **WHEN** Webhook 目标返回 HTTP 状态码 >= 400
- **THEN** 该 Webhook 计为失败

---

### Requirement: Localized email notification sent on download
`share:notify` Worker 任务 SHALL 对 `notify_emails` 中每个邮件地址通过 SMTP 发送通知邮件，邮件语言使用 `RedisShareInfo.Locale` 对应的翻译，不支持的 locale 回退到英文。

#### Scenario: Email sent in creator's locale
- **WHEN** SMTP 配置完整（`smtp.host` 非空）且 `shareInfo.Locale="zh-CN"`
- **THEN** Worker 向目标地址发送邮件，邮件 Subject 和 Body 为中文内容

#### Scenario: Email locale fallback to English
- **WHEN** `shareInfo.Locale` 为空或不支持的语言代码
- **THEN** Worker 使用英文模板发送邮件

#### Scenario: SMTP not configured
- **WHEN** `smtp.host` 为空
- **THEN** 邮件通知被跳过（不计入失败数），记录 warn 日志

---

### Requirement: Partial success means task success; all-fail triggers retry
`share:notify` 任务 SHALL 仅在所有通知目标（邮件 + Webhook）均失败时才返回 error；只要有一个成功，任务 SHALL 返回 nil。

#### Scenario: One webhook succeeds, one email fails
- **WHEN** 2 个通知目标中，Webhook 成功、Email 失败
- **THEN** 任务返回 nil（不重试）

#### Scenario: All notifications fail
- **WHEN** 所有 Webhook 和 Email 均失败
- **THEN** 任务返回聚合 error，Asynq 按默认策略重试

#### Scenario: Share expired when task executes
- **WHEN** Worker 执行时 Redis 中已无对应 `shareInfo`（share 已过期）
- **THEN** 任务返回 nil（静默跳过，不重试）
