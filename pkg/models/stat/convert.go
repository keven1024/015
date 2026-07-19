package stat

import "encoding/json"

func JsonStatDataToDomain(data string) (*StatData, error) {
	var stat StatData
	if err := json.Unmarshal([]byte(data), &stat); err != nil {
		return nil, err
	}
	return &stat, nil
}

func DomainStatDataToJson(stat *StatData) (string, error) {
	jsonData, err := json.Marshal(stat)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
