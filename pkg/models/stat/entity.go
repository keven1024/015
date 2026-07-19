package stat

// StatData 统计数据结构
type StatData struct {
	FileSize    int64 `json:"file_size"`    // 文件大小
	FileNum     int64 `json:"file_num"`     // 文件数量
	ShareNum    int64 `json:"share_num"`    // 分享数量
	DownloadNum int64 `json:"download_num"` // 下载数量
}
