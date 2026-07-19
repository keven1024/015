package task

import "encoding/json"

func JsonTaskInfoToDomain(data string) (*map[string]any, error) {
	var taskInfo map[string]any
	if err := json.Unmarshal([]byte(data), &taskInfo); err != nil {
		return nil, err
	}
	return &taskInfo, nil
}

func DomainTaskInfoToJson(taskInfo map[string]any) (string, error) {
	jsonData, err := json.Marshal(taskInfo)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
