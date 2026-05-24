## 1. 修复 pkg/i18n

- [x] 1.1 修改 `pkg/i18n/go.mod`：添加 `github.com/BurntSushi/toml`、`github.com/nicksnyder/go-i18n/v2`、`golang.org/x/text` 依赖
- [x] 1.2 修改 `pkg/i18n/i18n.go`：将 `Init()` 内部的 JSON 解析器替换为 TOML 解析器，文件后缀过滤从 `.json` 改为 `.toml`
- [x] 1.3 将 `worker/internal/i18n/` 下的所有 `.toml` 文件复制到 `pkg/i18n/locales/`（保留已有的 active.en.toml，添加 de/fr/ja/ko/zh-CN/zh-TW）
- [x] 1.4 在 `pkg/i18n/` 目录下运行 `go mod tidy`，生成 go.sum

## 2. 更新 go.work

- [x] 2.1 在 `go.work` 的 `use` 块中添加 `./pkg/i18n`
- [x] 2.2 在项目根目录运行 `go work sync` 更新 go.work.sum

## 3. 更新 worker

- [x] 3.1 在 `worker/go.mod` 中添加 `pkg/i18n` 依赖
- [x] 3.2 在 `worker/main.go` 中调用 `pkgi18n.Init()`，处理返回的 error
- [x] 3.3 修改 `worker/internal/tasks/notify.go`：删除 `loadI18nBundle`、`mustLocalize`、`localizeEmail` 三个函数及其相关 import，改用 `pkgi18n.TWithData` 生成邮件主题和正文
- [x] 3.4 在 `worker/` 目录下运行 `go mod tidy`，确认 `go-i18n`、`BurntSushi/toml`、`golang.org/x/text` 被正确移至 indirect 或移除

## 4. 清理

- [x] 4.1 删除 `worker/internal/i18n/` 目录
- [x] 4.2 在项目根目录运行 `go build ./...` 确认编译通过
