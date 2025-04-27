package models

type FileInfo struct {
	FileSize int64  `json:"size"`
	MimeType string `json:"mime_type"`
	FileHash string `json:"hash"`
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
	Expire    int64    `json:"expire"`
}

type RedisShareInfo struct {
	Id        string `json:"id"`
	Owner     string `json:"owner"`
	FileId    string `json:"fileId"`
	CreatedAt int64  `json:"created_at"`
}
