package models

import (
	"context"
	"encoding/json"
	"pkg/utils"
	"time"

	"github.com/redis/rueidis"
	"github.com/spf13/cast"
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
	UpdatedAt int64    `json:"updated_at"`
	Expire    int64    `json:"expire"` // 只有上传文件(init)的时候有这个字段
}

func GetRedisFileInfo(fileId string) (*RedisFileInfo, error) {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
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

func SetRedisFileInfo(fileId string, handler func(fileInfo *RedisFileInfo) *RedisFileInfo) (*RedisFileInfo, error) {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	old_fileInfo, err := GetRedisFileInfo(fileId)
	if err != nil {
		return nil, err
	}
	if old_fileInfo == nil {
		old_fileInfo = &RedisFileInfo{
			CreatedAt: time.Now().Unix(),
			Expire:    cast.ToInt64(utils.GetEnvWithDefault("upload.remove_expire", "2")) * 3600,
		}
	}
	fileInfo := handler(old_fileInfo)
	fileInfo.UpdatedAt = time.Now().Unix()
	jsonData, err := json.Marshal(fileInfo)
	if err != nil {
		return nil, err
	}
	if err := rdb.Do(ctx, rdb.B().Hset().Key("015:fileInfoMap").FieldValue().FieldValue(fileId, string(jsonData)).Build()).Error(); err != nil {
		return nil, err
	}
	return fileInfo, nil
}

func GetRedisFileInfoAll() (map[string]string, error) {
	rdb := utils.GetRedisClient()
	ctx := context.Background()
	return rdb.Do(ctx, rdb.B().Hgetall().Key("015:fileInfoMap").Build()).AsStrMap()
}
