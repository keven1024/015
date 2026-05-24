## Context

项目使用 Go workspace（go.work），`pkg/` 下有多个共享包（models、utils、services、i18n）。worker 模块通过 workspace 引用这些包。

当前 `pkg/i18n` 存在两个问题：
1. `Init()` 中注册的是 JSON 格式解析器，但 locales 目录下的文件是 `.toml`，导致实际加载翻译时什么都加载不到
2. 只有英文一种翻译，缺少其他语言

worker 独立维护了 7 种语言的翻译文件和一套完整的 i18n 初始化/加载逻辑，但这套逻辑使用动态路径查找（尝试多个相对路径），在容器化部署中脆弱。

## Goals / Non-Goals

**Goals:**
- 修复 `pkg/i18n` 的 TOML 加载 bug
- 将所有翻译文件集中到 `pkg/i18n/locales/`，通过 `embed.FS` 编译进二进制，消除运行时路径依赖
- worker 移除本地 i18n 实现，改用 `pkg/i18n`

**Non-Goals:**
- 不增加新的翻译 key
- 不修改 `Init()` 函数签名
- 不影响任何 HTTP API 或前端行为

## Decisions

### 1. 使用 embed.FS 而非动态路径

worker 当前通过 `filepath.Glob` 查找翻译文件，需要尝试 4 个不同路径。`pkg/i18n` 已使用 `//go:embed` 将文件编译进二进制，更可靠。翻译文件全部迁移到 `pkg/i18n/locales/`，worker 不再携带翻译文件。

### 2. 保留 Init() 函数不变

`Init()` 函数对外保持现有签名，只修改内部实现：将 `bundle.RegisterUnmarshalFunc("json", json.Unmarshal)` 改为 `bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)`，并将文件后缀过滤从 `.json` 改为 `.toml`。

### 3. worker 调用方式

worker 在 `main.go` 中调用 `pkgi18n.Init()`，`notify.go` 中将原来的三个本地函数替换为两次 `pkgi18n.TWithData()` 调用。

### 4. go.mod 依赖调整

`pkg/i18n/go.mod` 新增：
- `github.com/BurntSushi/toml`
- `github.com/nicksnyder/go-i18n/v2`
- `golang.org/x/text`

worker/go.mod 新增 `pkg/i18n`，如果 `go-i18n`、`BurntSushi/toml`、`golang.org/x/text` 不再被 worker 直接使用则移至 indirect。

## Risks / Trade-offs

- [翻译文件删除] 删除 `worker/internal/i18n/` 不可逆 → 迁移前确认 pkg/i18n/locales/ 中所有文件正确
- [go.mod 变动] worker 依赖关系变化可能引起 `go mod tidy` 意外移除仍需要的包 → 执行后检查 go.sum 和编译结果
