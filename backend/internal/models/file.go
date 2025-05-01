package models

import (
	"backend/internal/utils"
	"encoding/json"

	"dario.cat/mergo"
	"github.com/redis/go-redis/v9"
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
	Expire    int64    `json:"expire"`
}

type RedisShareInfo struct {
	Id        string `json:"id"`
	Owner     string `json:"owner"`
	FileId    string `json:"fileId"`
	CreatedAt int64  `json:"created_at"`
}

func GetRedisFileInfo(fileId string) (*RedisFileInfo, error) {
	rdb, ctx := utils.GetRedisClient()
	fileInfoUnmarshalData, err := rdb.HGet(ctx, "015:fileInfoMap", fileId).Result()
	if err == redis.Nil {
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
	_, err = rdb.HSet(ctx, "015:fileInfoMap", fileId, string(jsonData)).Result()
	return err
}
