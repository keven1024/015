package models

import (
	"encoding/json"
	"pkg/utils"

	"dario.cat/mergo"
	"github.com/redis/rueidis"
)

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
	Expire    int64    `json:"expire"` // 只有上传文件(init)的时候有这个字段
}

func GetRedisFileInfo(fileId string) (*RedisFileInfo, error) {
	rdb, ctx := utils.GetRedisClient()
	fileInfoUnmarshalData, err := rdb.Do(ctx, rdb.B().Hget().Key("015:fileInfoMap").Field(fileId).Build()).ToString()
	if rueidis.IsRedisNil(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var fileInfoData RedisFileInfo
	if err := json.Unmarshal([]byte(fileInfoUnmarshalData), &fileInfoData); err != nil {
		return nil, err
	}
	return &fileInfoData, nil
}

func SetRedisFileInfo(fileId string, fileInfo RedisFileInfo) error {
	rdb, ctx := utils.GetRedisClient()
	old_fileInfo, err := GetRedisFileInfo(fileId)
	if err != nil {
		return err
	}
	if old_fileInfo != nil {
		mergo.Merge(&fileInfo, old_fileInfo)
	}
	jsonData, _ := json.Marshal(fileInfo)
	return rdb.Do(ctx, rdb.B().Hset().Key("015:fileInfoMap").FieldValue().FieldValue(fileId, string(jsonData)).Build()).Error()
}

func GetRedisFileInfoAll() (map[string]string, error) {
	rdb, ctx := utils.GetRedisClient()
	return rdb.Do(ctx, rdb.B().Hgetall().Key("015:fileInfoMap").Build()).AsStrMap()
}
