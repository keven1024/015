package stat

const modelName = "015:stat"

type StatZSetKey string

const (
	FileZSetKey     StatZSetKey = "file"
	ShareZSetKey    StatZSetKey = "share"
	DownloadZSetKey StatZSetKey = "download"
	WorkerZSetKey   StatZSetKey = "worker"
)

// FileStatData 文件 zset 成员数据。
type FileStatData struct {
	FileSize int64  `json:"file_size"`
	FileID   string `json:"file_id"`
}

// ShareStatData 分享 zset 成员数据。
type ShareStatData struct {
	ShareID string `json:"share_id"`
}

// DownloadStatData 下载 zset 成员数据。
type DownloadStatData struct {
	ShareID string `json:"share_id"`
	IP      string `json:"ip"`
}

// WorkerStatData worker zset 成员数据。
type WorkerStatData struct {
	TaskID   string `json:"task_id"`
	TaskType string `json:"task_type"`
	Status   string `json:"status"`
}

type RedisStatDataItem struct {
	Timestamp int64
	Data      string
}
