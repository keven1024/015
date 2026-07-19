package file_share_relational

import "encoding/json"

func JsonFileShareRelationalToDomain(data string) ([]string, error) {
	var shareIDs []string
	if err := json.Unmarshal([]byte(data), &shareIDs); err != nil {
		return nil, err
	}
	return shareIDs, nil
}

func DomainFileShareRelationalToJson(shareIDs []string) (string, error) {
	jsonData, err := json.Marshal(shareIDs)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
