package stat

import "encoding/json"

func JsonFileStatDataToDomain(data string) (*FileStatData, error) {
	var stat FileStatData
	if err := json.Unmarshal([]byte(data), &stat); err != nil {
		return nil, err
	}
	return &stat, nil
}

func DomainFileStatDataToJson(stat *FileStatData) (string, error) {
	jsonData, err := json.Marshal(stat)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func JsonShareStatDataToDomain(data string) (*ShareStatData, error) {
	var stat ShareStatData
	if err := json.Unmarshal([]byte(data), &stat); err != nil {
		return nil, err
	}
	return &stat, nil
}

func DomainShareStatDataToJson(stat *ShareStatData) (string, error) {
	jsonData, err := json.Marshal(stat)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func JsonDownloadStatDataToDomain(data string) (*DownloadStatData, error) {
	var stat DownloadStatData
	if err := json.Unmarshal([]byte(data), &stat); err != nil {
		return nil, err
	}
	return &stat, nil
}

func DomainDownloadStatDataToJson(stat *DownloadStatData) (string, error) {
	jsonData, err := json.Marshal(stat)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func JsonWorkerStatDataToDomain(data string) (*WorkerStatData, error) {
	var stat WorkerStatData
	if err := json.Unmarshal([]byte(data), &stat); err != nil {
		return nil, err
	}
	return &stat, nil
}

func DomainWorkerStatDataToJson(stat *WorkerStatData) (string, error) {
	jsonData, err := json.Marshal(stat)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
