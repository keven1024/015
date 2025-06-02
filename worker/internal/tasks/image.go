package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"worker/internal/models"
	"worker/internal/services"
	"worker/internal/utils"

	"github.com/hibiken/asynq"
)

type ImageTaskPayload struct {
	FileId string `json:"file_id"`
}

type CompressImageTaskPayload struct {
	ImageTaskPayload
}

func CompressImage(ctx context.Context, task *asynq.Task) error {
	var payload CompressImageTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	originalFileInfo, _ := models.GetRedisFileInfo(payload.FileId)
	if originalFileInfo == nil || originalFileInfo.FileType != models.FileTypeUpload {
		return errors.New("文件不存在")
	}
	uploadPath, err := utils.GetUploadDirPath()
	if err != nil {
		return err
	}
	originalPath := filepath.Join(uploadPath, payload.FileId)
	switch originalFileInfo.MimeType {
	case "image/png":
		args := []string{"--output", originalPath + "_compressed", originalPath}
		cmd := exec.Command("pngquant", args...)
		_, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
	case "image/jpeg":
		err := utils.CopyFile(originalPath, originalPath+"_compressed")
		if err != nil {
			return err
		}
		args := []string{"-m", "90", "--strip-all", originalPath + "_compressed"}
		cmd := exec.Command("jpegoptim", args...)
		_, err = cmd.CombinedOutput()
		if err != nil {
			return err
		}
	default:
		return errors.New("不支持的文件类型")
	}
	compressedPath := fmt.Sprintf("%s_compressed", originalPath)
	compressedFileInfo, err := services.GenStandardFile(compressedPath, originalFileInfo.MimeType)
	if err != nil {
		return err
	}

	models.SetRedisTaskInfo(task.ResultWriter().TaskID(), map[string]any{
		"status": "success",
		"result": []any{
			map[string]any{
				"old_file": map[string]any{
					"id":   payload.FileId,
					"size": originalFileInfo.FileSize,
				},
				"new_file": map[string]any{
					"id":   compressedFileInfo.FileId,
					"size": compressedFileInfo.FileSize,
				},
			},
		},
	})

	return nil
}

func UpscaleImage(ctx context.Context, task *asynq.Task) error {
	return nil
}

func TranslateImage(ctx context.Context, task *asynq.Task) error {
	return nil
}

func CreateAIImage(ctx context.Context, task *asynq.Task) error {
	return nil
}
