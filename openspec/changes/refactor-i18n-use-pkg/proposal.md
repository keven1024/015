## Why

worker 模块自己实现了一套独立的 i18n 加载逻辑（动态文件路径查找、手动注册 TOML 解析），而 `pkg/i18n` 已经提供了通用的 i18n 封装但从未被使用，且其实现存在 bug（注册的是 JSON 格式，实际文件是 TOML）。统一到 `pkg/i18n` 可以消除重复、修复 bug、让后续模块复用。

## What Changes

- 修复 `pkg/i18n`：将格式从 JSON 改为 TOML，更新 `go.mod` 添加缺失依赖
- 将 worker/internal/i18n/ 的 7 种语言翻译文件迁移到 `pkg/i18n/locales/`
- 将 `pkg/i18n` 加入 `go.work` 工作区
- worker/go.mod 中添加 `pkg/i18n` 依赖，移除对 `go-i18n`、`BurntSushi/toml`、`golang.org/x/text/language` 的直接引用
- 删除 worker 中的 `loadI18nBundle`、`mustLocalize`、`localizeEmail` 函数，改用 `pkg/i18n.TWithData`
- 删除 `worker/internal/i18n/` 目录

## Capabilities

### New Capabilities

无新能力，纯重构。

### Modified Capabilities

- `share-download-notify`：i18n 实现层迁移到 `pkg/i18n`，接口行为（邮件本地化主题和正文）不变

## Impact

- **pkg/i18n**：go.mod 新增 `github.com/BurntSushi/toml`、`github.com/nicksnyder/go-i18n/v2`、`golang.org/x/text` 依赖，locales 目录新增 6 种语言文件
- **go.work**：新增 `./pkg/i18n` 模块
- **worker**：notify.go 逻辑简化，go.mod 依赖变化，`worker/internal/i18n/` 目录删除
- 不影响任何 API 或前端
