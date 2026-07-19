package file

import "encoding/json"

func JsonFileInfoToDomain(data string) (*RedisFileInfo, error) {
	var fileInfo RedisFileInfo
	if err := json.Unmarshal([]byte(data), &fileInfo); err != nil {
		return nil, err
	}
	return &fileInfo, nil
}

func DomainFileInfoToJson(fileInfo *RedisFileInfo) (string, error) {
	jsonData, err := json.Marshal(fileInfo)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
