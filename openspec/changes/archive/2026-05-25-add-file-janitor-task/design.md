## Context

worker 目前所有檔案清理都依賴 `file:remove` 延遲任務，由業務邏輯在適當時機主動排程（上傳完成後設定 TTL、share 刪除時觸發）。若任務在排程前 worker 重啟、Redis 任務丟失或業務邏輯有 bug，檔案就會永久殘留在磁碟上。`asynq` 已提供 `Scheduler` 元件支援 cron 排程，無需引入新依賴。

## Goals / Non-Goals

**Goals:**
- 每天凌晨自動掃描並清理三類漏網檔案
- 不重複造輪子，清理邏輯複用現有 `file:remove` 任務
- 不影響現有任務流程與 API

**Non-Goals:**
- 不處理 `Expire == 0` 的特殊情況
- 不清理 `fileShareRelational` 裡的孤兒 share 條目
- 不提供手動觸發的 HTTP 介面

## Decisions

**1. 使用 `asynq.Scheduler` 而非外部 cron，排程時間為每天凌晨 03:00**

`asynq.Scheduler` 在同一 worker 程式內以獨立 goroutine 運行，透過 Redis 做 leader election（避免多實例重複觸發）。外部 cron 需要額外基礎設施。選擇 Scheduler，cron 表達式為 `0 3 * * *`。

```
main.go
  ├── asynq.NewServer(...)     // 處理任務
  ├── asynq.NewScheduler(...)  // 排程 file:janitor @ 03:00
  └── mux.HandleFunc("file:janitor", tasks.FileJanitor)
```

**2. `FileJanitor` 函式放在 `worker/internal/tasks/file.go`，超過 300 行才拆資料夾**

目前 `file.go` 為 50 行，加入 janitor 後預計約 100 行，遠低於 300 行門檻，直接寫入 `file.go`。若未來該檔超過 300 行，則改為 `worker/internal/tasks/file/` 目錄，拆成 `janitor.go`（清理邏輯）與 `remove.go`（原 `RemoveFile` 邏輯）。

**3. Case 1（孤兒本地檔）直接 `os.RemoveAll`，不走 `file:remove`**

這類檔案在 fileInfoMap 中不存在，`RemoveFile` handler 會直接 return nil（找不到 fileInfo），走 `file:remove` 無法刪除磁碟檔案。必須直接刪除。

**4. Case 2、3 透過 `SetFileRemoveTask(id, 0)` 委派**

複用 `RemoveFile` 現有邏輯（含 share 二次確認），避免重複實作。delay=0 表示立即處理。

**5. 掃描策略：一次性全量讀取**

- 本地目錄：`os.ReadDir(uploadPath)` 取得所有檔名（即 fileId）
- Redis：`GetRedisFileInfoAll()` 一次取得完整 fileInfoMap
- 兩者建立 map 後交叉比對，時間複雜度 O(n)

```
本地檔案 set  ────┐
                  ├─ 差集 → Case 1
Redis fileInfoMap ┘

Redis fileInfoMap ─ 遍歷 ─┬─ type==init && expired → Case 2
                           └─ type==already && no share → Case 3
```

**6. Case 3 的 share 關係查詢**

對每個 `type==already` 的 fileId 呼叫 `GetRedisFileShareRelational(fileId)`，若回傳空 slice 則觸發刪除。此為 O(n) Redis 查詢，數量級與檔案數相同，凌晨低峰期執行可接受。

## Risks / Trade-offs

- **誤刪風險（Case 1）**：本地檔案剛建立但 fileInfoMap 尚未寫入的時間窗口（上傳初始化中）。→ 凌晨執行，正在上傳中的檔案其 init 記錄通常已存在 Redis，風險極低。
- **Redis 查詢量（Case 3）**：若檔案數量龐大，逐一查詢 fileShareRelational 會產生大量 Redis 請求。→ 目前是輕量平台，可接受；未來可改用 `HGETALL fileShareRelational` 一次取得再做本地比對。
- **Scheduler 單點**：`asynq.Scheduler` 依賴 Redis leader election，Redis 不可用時排程失敗。→ 可接受，Redis 不可用時整個 worker 本就無法運作。

## Migration Plan

直接部署，無資料遷移。首次執行會清理歷史殘留檔案，屬於預期行為。
