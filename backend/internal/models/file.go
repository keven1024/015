package models

import (
	"backend/internal/utils"
	"encoding/json"
	"errors"

	"dario.cat/mergo"
)

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

func GetRedisFileInfo(fileId string) (RedisFileInfo, error) {
	rdb, ctx := utils.GetRedisClient()
	fileInfoUnmarshalData, _ := rdb.HGet(ctx, "015:fileInfoMap", fileId).Result()

	if fileInfoUnmarshalData != "" {
		var fileInfoData RedisFileInfo
		if err := json.Unmarshal([]byte(fileInfoUnmarshalData), &fileInfoData); err != nil {
			return RedisFileInfo{}, err
		}
		return fileInfoData, nil
	}
	return RedisFileInfo{}, errors.New("db不存在该文件信息")
}

func SetRedisFileInfo(fileId string, fileInfo RedisFileInfo) error {
	rdb, ctx := utils.GetRedisClient()
	old_fileInfo, err := GetRedisFileInfo(fileId)
	if err != nil {
		return err
	}
	mergo.Merge(&fileInfo, old_fileInfo)
	jsonData, _ := json.Marshal(fileInfo)
	_, err = rdb.HSet(ctx, "015:fileInfoMap", fileId, string(jsonData)).Result()
	return err
}
