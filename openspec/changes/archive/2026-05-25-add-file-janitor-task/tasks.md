## 1. 新增 FileJanitor 任務處理函式

- [x] 1.1 在 `worker/internal/tasks/file.go` 末尾新增 `FileJanitor(ctx context.Context, task *asynq.Task) error`（若加入後超過 300 行，則改建 `worker/internal/tasks/file/` 目錄，將 `RemoveFile` 移至 `remove.go`，`FileJanitor` 放至 `janitor.go`）
- [x] 1.2 在 `FileJanitor` 中呼叫 `u.GetUploadDirPath()` 取得本地上傳目錄路徑，並用 `os.ReadDir` 掃描目錄取得本地檔案清單
- [x] 1.3 呼叫 `models.GetRedisFileInfoAll()` 取得 fileInfoMap 全量資料，建立 fileId set
- [x] 1.4 實作 Case 1：本地存在但 fileInfoMap 無記錄的 fileId，呼叫 `os.RemoveAll(filePath)` 直接刪除
- [x] 1.5 遍歷 fileInfoMap，實作 Case 2：`FileType == "init"` 且 `CreatedAt + Expire < time.Now().Unix()`，呼叫 `pkgservices.SetFileRemoveTask(fileId, 0)`
- [x] 1.6 遍歷 fileInfoMap，實作 Case 3：`FileType == "already"` 且 `GetRedisFileShareRelational(fileId)` 回傳空 slice，呼叫 `pkgservices.SetFileRemoveTask(fileId, 0)`

## 2. 在 worker/main.go 中整合

- [x] 2.1 在 `mux` 中註冊 `mux.HandleFunc("file:janitor", tasks.FileJanitor)`
- [x] 2.2 建立 `asynq.NewScheduler`，使用與 `asynq.NewServer` 相同的 Redis 連線設定
- [x] 2.3 呼叫 `scheduler.Register("0 3 * * *", asynq.NewTask("file:janitor", nil))` 設定每日凌晨 03:00 排程
- [x] 2.4 在 `srv.Run(mux)` 之前啟動 `scheduler.Start()`，並在程式結束時呼叫 `scheduler.Shutdown()`

## 3. 驗證

- [x] 3.1 在 `worker/` 目錄下執行 `go build ./...`，確認編譯通過
