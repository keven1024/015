package share

import "encoding/json"

func JsonShareInfoToDomain(data string) (*RedisShareInfo, error) {
	var shareInfo RedisShareInfo
	if err := json.Unmarshal([]byte(data), &shareInfo); err != nil {
		return nil, err
	}

	// 兼容旧data逻辑，后续删除
	payload := map[string]json.RawMessage{}
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		return nil, err
	}
	if rawData, ok := payload["data"]; ok {
		var legacyData string
		if err := json.Unmarshal(rawData, &legacyData); err != nil {
			return nil, err
		}
		if shareInfo.Type == ShareTypeFile && len(shareInfo.Files) == 0 {
			files, err := JsonShareFilesToDomain(legacyData)
			if err != nil {
				return nil, err
			}
			shareInfo.Files = files
		}
		if shareInfo.Type == ShareTypeText && shareInfo.Text == "" {
			shareInfo.Text = legacyData
		}
	}
	return &shareInfo, nil
}

func DomainShareInfoToJson(shareInfo *RedisShareInfo) (string, error) {
	jsonData, err := json.Marshal(shareInfo)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func JsonShareFilesToDomain(data string) ([]ShareFileData, error) {
	if data == "" {
		return nil, nil
	}
	var files []ShareFileData
	if err := json.Unmarshal([]byte(data), &files); err == nil {
		return files, nil
	}
	return []ShareFileData{{
		Id:       data,
		FileName: data,
	}}, nil
}
