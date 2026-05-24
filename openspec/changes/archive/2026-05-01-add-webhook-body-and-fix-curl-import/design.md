## 背景

后端（`pkg/models/share.go` 的 `NotifyWebhook`、`worker/internal/tasks/notify.go` 的 `sendWebhook`）已完整支持 `Body` 和 `BodyType`。前端 `WebhookItem` 接口和表单 UI 缺少 `body` 字段，`NotifyConfigField.vue` 中的 curl 导入处理函数仅为存根（解析命令后只 `console.log` 结果，未写入表单）。

前端 headers 存储格式为 `[string, string][]`（键值对数组），与 `KvInputGroupField` 的输出一致。后端模型使用 `map[string]string`——此不一致为既有问题，不在本次变更范围内。

## 目标 / 非目标

**目标：**
- 在 Webhook 高级面板中暴露 `body` textarea，供用户配置请求体。
- 激活 curl 导入处理：粘贴 `curl` 命令 → 自动填入 URL、method、headers、body，并展开高级面板。
- 为新字段添加 i18n 键。

**非目标：**
- `bodyType` 选择器（form-data vs raw）——body 默认以原始文本处理，本次不暴露 `bodyType` 字段。
- 修复前端 `[string, string][]` 与后端 `map[string]string` 的既有格式不一致问题。
- 后端及 Worker 变更——它们已支持 body。

## 决策

**D1 — 不在 UI 中暴露 `bodyType` 字段**
Worker 已对 `bodyType` 做分支处理（form-data / raw / none），但增加类型选择器会提升 UI 复杂度。默认使用 raw（`bodyType` 未设置或为空 → body 原样发送）覆盖了大多数场景，也符合用户粘贴 curl 命令时的预期。后续可单独追加，不会产生 breaking change。

**D2 — curl 导入后自动展开高级面板**
curl 解析成功后，headers 和 body 已被填入，但高级面板处于折叠状态，用户不可见。自动展开可让用户即时查看并验证填入的数据，避免出现隐式状态。

**D3 — curl 导入只预填表单，不自动提交**
导入仅预填表单字段，用户仍需手动点击保存按钮，避免意外副作用。

## 风险 / 权衡

- **`parseCurl` 的 `body` 字段不一定存在** — 对于 GET 请求，`sweet-curl-parser` 可能不填充 `data.body`。处理函数需用 `if (data.body)` 做判断再写入。→ 缓解：条件赋值，新建 Webhook 时默认 body 为空字符串。
- **curl headers 格式转换** — `parseCurl` 返回 `{name: string, value: string}[]`，表单期望 `[string, string][]`，转换逻辑为 `data.headers.map(h => [h.name, h.value])`。→ 无风险，但须显式处理。
