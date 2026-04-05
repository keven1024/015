package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"pkg/models"
	"pkg/utils"
	"worker/internal/services"

	"github.com/hibiken/asynq"
)

func CompressImage(ctx context.Context, task *asynq.Task) error {
	var payload CompressImageTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	originalFileInfo, _ := models.GetRedisFileInfo(payload.FileId)
	if originalFileInfo == nil || originalFileInfo.FileType != models.FileTypeUpload {
		return ErrNotFoundFile
	}
	uploadPath, err := utils.GetUploadDirPath()
	if err != nil {
		return err
	}
	originalPath := filepath.Join(uploadPath, payload.FileId)
	compressedPath, err := services.CompressImage(originalPath, originalFileInfo.MimeType)
	if err != nil {
		if errors.Is(err, services.ErrUnsupportedMimeType) {
			return fmt.Errorf("%w: %w", err, asynq.SkipRetry)
		}
		return err
	}
	compressedFileInfo, err := services.GenStandardFile(compressedPath, originalFileInfo.MimeType)
	if err != nil {
		defer os.Remove(compressedPath) //nolint:errcheck
		return err
	}

	if err := models.SetRedisTaskInfo(task.ResultWriter().TaskID(), map[string]any{
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
	}); err != nil {
		return err
	}

	return nil
}

func ConvertImage(ctx context.Context, task *asynq.Task) error {
	var payload ConvertImageTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	originalFileInfo, _ := models.GetRedisFileInfo(payload.FileId)
	if originalFileInfo == nil || originalFileInfo.FileType != models.FileTypeUpload {
		return ErrNotFoundFile
	}
	uploadPath, err := utils.GetUploadDirPath()
	if err != nil {
		return err
	}
	originalPath := filepath.Join(uploadPath, payload.FileId)
	convertedPath, err := services.ConvertImageWithMagick(originalPath, originalFileInfo.MimeType, payload.TargetExt)
	if err != nil {
		if errors.Is(err, services.ErrUnsupportedMimeType) {
			return fmt.Errorf("%w: %w", err, asynq.SkipRetry)
		}
		return err
	}
	mimeType := mime.TypeByExtension(fmt.Sprintf(".%s", payload.TargetExt))
	if mimeType == "" {
		defer os.Remove(convertedPath) //nolint:errcheck
		return ErrUnknown
	}
	convertedFileInfo, err := services.GenStandardFile(convertedPath, mimeType)
	if err != nil {
		defer os.Remove(convertedPath) //nolint:errcheck
		return err
	}

	if err := models.SetRedisTaskInfo(task.ResultWriter().TaskID(), map[string]any{
		"status": "success",
		"result": []any{
			map[string]any{
				"old_file": map[string]any{
					"id":   payload.FileId,
					"size": originalFileInfo.FileSize,
				},
				"new_file": map[string]any{
					"id":   convertedFileInfo.FileId,
					"size": convertedFileInfo.FileSize,
				},
			},
		},
	}); err != nil {
		return err
	}

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
