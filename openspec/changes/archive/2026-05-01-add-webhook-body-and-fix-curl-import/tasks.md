## 1. 前端 — 类型定义与默认值

- [x] 1.1 在 `NotifyConfigField.vue` 的 `WebhookItem` 接口中新增 `body?: string` 字段
- [x] 1.2 将"添加"按钮的默认值从 `{ url: '', method: 'POST', headers: [] }` 改为 `{ url: '', method: 'POST', headers: [], body: '' }`

## 2. 前端 — Body Textarea UI

- [x] 2.1 在 `NotifyConfigField.vue` 中从 `'../Field/TextareaField.vue'` 导入 `TextareaField`
- [x] 2.2 在 `v-show="expandedAdvanced.has(index)"` 区块内的 `<KvInputField>` 下方添加 `<TextareaField>`，绑定到 `notify_webhooks.${index}.body`，label 使用 i18n 键 `page.shareOptions.notify.webhookBody`，`rows=4`，placeholder 为 `{"key": "value"}`

## 3. 前端 — curl 导入处理函数

- [x] 3.1 将 `console.log('command', data)` 及注释掉的代码块替换为实际赋值逻辑：
  - `setFieldValue(`notify_webhooks.${index}.url`, data.url.fullUrl)`
  - `setFieldValue(`notify_webhooks.${index}.method`, data.method.toUpperCase())`
  - `setFieldValue(`notify_webhooks.${index}.headers`, data.headers.map((h: any) => [h.name, h.value]))`
  - `if (data.body) setFieldValue(`notify_webhooks.${index}.body`, data.body)`
- [x] 3.2 curl 解析成功后，将当前 index 加入 `expandedAdvanced` 以自动展开高级面板：`expandedAdvanced = new Set([...expandedAdvanced, index])`

## 4. i18n

- [x] 4.1 在 `front/i18n/locales/en.json` 的 `"webhookHeaders"` 键后新增 `"webhookBody": "Request Body"`
- [x] 4.2 在 `front/i18n/locales/zh-CN.json` 中新增 `"webhookBody": "请求体"`
- [x] 4.3 在 `front/i18n/locales/zh-TW.json` 中新增 `"webhookBody": "請求體"`
- [x] 4.4 在 `front/i18n/locales/de.json` 中新增 `"webhookBody": "Anfrage-Body"`
- [x] 4.5 在 `front/i18n/locales/fr.json` 中新增 `"webhookBody": "Corps de la requête"`
- [x] 4.6 在 `front/i18n/locales/ja.json` 中新增 `"webhookBody": "リクエストボディ"`
- [x] 4.7 在 `front/i18n/locales/ko.json` 中新增 `"webhookBody": "요청 본문"`
