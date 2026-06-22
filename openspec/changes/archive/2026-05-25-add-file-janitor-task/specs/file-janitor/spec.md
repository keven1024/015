## ADDED Requirements

### Requirement: 每日定期排程 file:janitor 任務
系統 SHALL 在每天凌晨 00:00（伺服器本地時間）自動排程一次 `file:janitor` asynq 任務。多個 worker 實例並存時，系統 SHALL 保證同一時間只有一個實例觸發排程（透過 asynq.Scheduler 的 Redis leader election 機制）。

#### Scenario: 單實例正常排程
- **WHEN** worker 啟動且時間到達每天 00:00
- **THEN** asynq.Scheduler 將 `file:janitor` 任務加入 asynq 佇列一次

#### Scenario: 多實例防重複觸發
- **WHEN** 多個 worker 實例同時運行且時間到達 00:00
- **THEN** 只有一個實例成功排程，其餘實例因 leader election 未獲鎖而跳過

### Requirement: 清理本地孤兒檔案（無 fileInfoMap 記錄）
系統 SHALL 掃描本地上傳目錄，對每個存在於磁碟但在 Redis `fileInfoMap` 中無對應記錄的檔案，直接執行 `os.RemoveAll` 刪除。

#### Scenario: 刪除孤兒本地檔案
- **WHEN** 本地上傳目錄中存在 fileId 目錄，且 Redis fileInfoMap 中無該 fileId 的記錄
- **THEN** 系統直接刪除該本地目錄

#### Scenario: 保留有 fileInfoMap 記錄的檔案
- **WHEN** 本地上傳目錄中存在 fileId 目錄，且 Redis fileInfoMap 中有該 fileId 的記錄
- **THEN** 系統不刪除該本地目錄，繼續處理下一個

### Requirement: 清理已過期的 init 狀態檔案
系統 SHALL 遍歷 Redis `fileInfoMap`，對 `FileType == "init"` 且滿足 `CreatedAt + Expire < 當前 Unix 時間戳` 的記錄，透過 `SetFileRemoveTask(fileId, 0)` 排程立即刪除。

#### Scenario: 排程刪除已過期 init 檔案
- **WHEN** fileInfoMap 中存在 FileType=="init" 且 CreatedAt+Expire < now() 的記錄
- **THEN** 系統呼叫 SetFileRemoveTask(fileId, 0)，將 file:remove 任務加入佇列

#### Scenario: 保留未過期的 init 檔案
- **WHEN** fileInfoMap 中存在 FileType=="init" 且 CreatedAt+Expire >= now() 的記錄
- **THEN** 系統不排程刪除，繼續處理下一個

### Requirement: 清理無 share 關係的已完成上傳檔案
系統 SHALL 遍歷 Redis `fileInfoMap`，對 `FileType == "already"` 且在 Redis `fileShareRelational` 中無任何 share 關係的記錄，透過 `SetFileRemoveTask(fileId, 0)` 排程立即刪除。

#### Scenario: 排程刪除無 share 關係的已完成檔案
- **WHEN** fileInfoMap 中存在 FileType=="already" 的記錄，且 fileShareRelational 中該 fileId 對應的 shareId 列表為空或不存在
- **THEN** 系統呼叫 SetFileRemoveTask(fileId, 0)，將 file:remove 任務加入佇列

#### Scenario: 保留有 share 關係的已完成檔案
- **WHEN** fileInfoMap 中存在 FileType=="already" 的記錄，且 fileShareRelational 中該 fileId 有一個或以上 shareId
- **THEN** 系統不排程刪除，繼續處理下一個
