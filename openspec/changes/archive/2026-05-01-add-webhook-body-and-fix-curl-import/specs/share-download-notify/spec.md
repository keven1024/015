## MODIFIED Requirements

### Requirement: Webhook notification sent on download
`share:notify` Worker 任务 SHALL 对 `notify_webhooks` 中每个配置项发起 HTTP 请求。

#### Scenario: POST webhook with form-data body
- **WHEN** Webhook `method=POST`、`bodyType=form-data`、`body="k=v"`
- **THEN** Worker 向 `url` 发起 POST 请求，`Content-Type: application/x-www-form-urlencoded`，Body 为 `k=v`

#### Scenario: POST webhook with raw body
- **WHEN** Webhook `method=POST`、`bodyType=raw` 或 `bodyType` 为空、`body` 非空
- **THEN** Worker 向 `url` 发起 POST 请求，Body 原样发送，不强制设置 Content-Type

#### Scenario: Webhook body empty or absent
- **WHEN** `body` 为空字符串或未设置
- **THEN** Worker 发起请求时不附加 Body

#### Scenario: Webhook with custom headers
- **WHEN** Webhook `headers` 中包含自定义请求头
- **THEN** Worker 将这些 Header 附加到请求中

#### Scenario: Webhook returns 4xx/5xx
- **WHEN** Webhook 目标返回 HTTP 状态码 >= 400
- **THEN** 该 Webhook 计为失败
