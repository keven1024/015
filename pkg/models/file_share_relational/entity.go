package file_share_relational

const modelName = "015:fileShareRelational"

func RedisLockKey(fileId string) string {
	return modelName + ":" + fileId
}
