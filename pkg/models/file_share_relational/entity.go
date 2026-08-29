package file_share_relational

func RedisLockKey(fileId string) string {
	return ModelName + ":" + fileId
}
