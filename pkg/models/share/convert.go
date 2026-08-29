package share

import (
	"encoding/json"
)

func JsonShareInfoToDomain(data string) (*RedisShareInfo, error) {
	var shareInfo RedisShareInfo
	if err := json.Unmarshal([]byte(data), &shareInfo); err != nil {
		return nil, err
	}
	// 旧数据兼容
	if shareInfo.Data != "" {
		switch shareInfo.Type {
		case ShareTypeFile:
			if len(shareInfo.Files) == 0 {
				shareInfo.Files = []ShareFileData{{
					Id:       shareInfo.Data,
					FileName: shareInfo.FileName,
				}}
			}
		case ShareTypeText:
			if shareInfo.Text == "" {
				shareInfo.Text = shareInfo.Data
			}
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
