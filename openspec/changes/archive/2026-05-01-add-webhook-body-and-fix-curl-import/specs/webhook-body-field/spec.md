## ADDED Requirements

### Requirement: Webhook body 字段可在 UI 中配置
系统 SHALL 在 Webhook 高级设置面板中暴露一个 `body` textarea，允许用户为每条 Webhook 条目输入原始请求体字符串。

#### Scenario: 用户输入 body 文本
- **WHEN** 用户展开某条 Webhook 的高级面板并在 body textarea 中输入文本
- **THEN** body 值以 `notify_webhooks[n].body` 的路径存储到表单中

#### Scenario: body 字段默认为空
- **WHEN** 用户通过"添加"按钮新增一条 Webhook 条目
- **THEN** `body` 默认值为空字符串

---

### Requirement: curl 导入自动填充 method、headers 和 body
系统 SHALL 解析粘贴到 Webhook URL 字段的 `curl` 命令，并自动将解析结果中的 `url`、`method`、`headers`、`body` 填入表单，同时展开高级面板。

#### Scenario: 将包含 headers 的有效 curl 命令粘贴到 URL 字段
- **WHEN** 用户在 URL 输入框中粘贴以 `curl ` 开头的字符串，且字段失去焦点
- **THEN** `url` 被设置为 `data.url.fullUrl`，`method` 被设置为 `data.method`（大写），`headers` 被转换为 `[string, string][]` 键值对数组，对应 Webhook 条目的高级面板被展开

#### Scenario: curl 命令包含 body
- **WHEN** 解析后的 curl 结果中包含非空的 `body` 字段
- **THEN** `body` 被写入 `notify_webhooks[n].body`

#### Scenario: curl 命令不包含 body（如 GET 请求）
- **WHEN** 解析后的 curl 结果中没有 `body` 字段或其值为空
- **THEN** `notify_webhooks[n].body` 保持为空字符串（不被 undefined 覆盖）

#### Scenario: 解析失败或输入不是 curl 命令
- **WHEN** URL 字段值不以 `curl ` 开头，或 `parseCurl` 返回 `success: false`
- **THEN** 表单中所有字段保持不变
