package services

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"time"

	"pkg/utils"

	"github.com/hibiken/asynq"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func SetFileRemoveTask(fileId string, expire time.Duration) error {
	client, payload, err := getCreateTaskMetadata(fileId)
	if err != nil {
		return err
	}
	_, err = client.Enqueue(
		asynq.NewTask("file:remove", payload),
		asynq.ProcessIn(expire),
		asynq.TaskID(fileId),
	)
	return err
}

func SetUploadRemoveTask(fileId string, expire time.Duration) error {
	client, payload, err := getCreateTaskMetadata(fileId)
	if err != nil {
		return err
	}
	_, err = client.Enqueue(
		asynq.NewTask("upload:remove", payload),
		asynq.ProcessIn(expire),
	)
	return err
}

func GetFileSliceList(fileId string, root *os.Root) ([]int, error) {
	files, err := fs.ReadDir(root.FS(), fmt.Sprintf("%s_%s", fileId, "tmp"))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %w", ErrReadFileSliceDir, err)
	}
	fileSliceList := lo.FilterMap(files, func(file os.DirEntry, _ int) (int, bool) {
		index, err := cast.ToIntE(file.Name())
		return index, err == nil
	})
	sort.Ints(fileSliceList)
	return fileSliceList, nil
}

func getCreateTaskMetadata(fileId string) (*asynq.Client, []byte, error) {
	client := utils.GetQueueClient()
	payload, err := json.Marshal(map[string]any{
		"file_id": fileId,
	})
	if err != nil {
		return nil, nil, err
	}
	return client, payload, nil
}
