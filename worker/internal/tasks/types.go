package tasks

import "errors"

var (
	ErrNotFoundFile = errors.New("NotFoundFile")
	ErrUnknown      = errors.New("Unknown")
)

type BaseFileTaskPayload struct {
	FileId string `json:"file_id"`
}

type RemoveFileTaskPayload struct {
	BaseFileTaskPayload
}

type ShareNotifyTaskPayload struct {
	ShareId string `json:"share_id"`
	IP      string `json:"ip"`
}

type CompressImageTaskPayload struct {
	BaseFileTaskPayload
}

type ConvertImageTaskPayload struct {
	BaseFileTaskPayload
	TargetExt string `json:"target_ext"`
}

type TranslateTextTaskPayload struct {
	Text     string `json:"text"`
	Source   string `json:"source"`
	Target   string `json:"target"`
	Provider string `json:"provider"`
}
