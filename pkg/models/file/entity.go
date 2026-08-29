package file

const modelName = "015:fileInfoMap"

func redisLockKey(fileId string) string {
	return modelName + ":" + fileId
}

type FileInfo struct {
	FileSize  int64  `json:"size"`
	MimeType  string `json:"mime_type"`
	FileHash  string `json:"hash"`
	ChunkSize int64  `json:"chunk_size"`
}

type FileType string

const (
	FileTypeInit   FileType = "init"
	FileTypeUpload FileType = "already"
)

type RedisFileInfo struct {
	FileInfo
	FileType  FileType `json:"type"`
	CreatedAt int64    `json:"created_at"`
	UpdatedAt int64    `json:"updated_at"`
	Expire    int64    `json:"expire"` // 只有上传文件(init)的时候有这个字段
}
