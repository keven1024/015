## Why

現有的檔案清理機制依賴 asynq 延遲任務（`file:remove`），在異常情況下（worker 重啟、任務丟失、Redis 資料不一致）可能導致本地磁碟殘留孤兒檔案，長期累積佔用儲存空間。需要一個每日定期掃描的兜底清理機制。

## What Changes

- 新增 `file:janitor` asynq 任務處理函式（`worker/internal/tasks/janitor.go`）
- 在 `worker/main.go` 中加入 `asynq.Scheduler`，每天凌晨 00:00 自動排程 `file:janitor`
- 在 `worker/main.go` 的 mux 中註冊 `file:janitor` handler

## Capabilities

### New Capabilities

- `file-janitor`: 每日定期掃描並清理三類漏網檔案：本地有但 fileInfoMap 無的孤兒檔、init 狀態已過期的未完成上傳、已完成上傳但無任何 share 關係的孤立檔

### Modified Capabilities

## Impact

- **worker/internal/tasks/**：新增 `janitor.go`
- **worker/main.go**：新增 `asynq.Scheduler` 及 `file:janitor` handler 註冊
- 不影響任何 API 或前端
- 無新增外部依賴（`asynq.Scheduler` 已在現有依賴中）
