package utils

import "go.uber.org/zap"

var sugar *zap.SugaredLogger

func init() {
	InitLogClient()
}

func InitLogClient() {
	if sugar != nil {
		return
	}
	logger, _ := zap.NewProduction()
	sugar = logger.Sugar()
}

func GetLogClient() *zap.SugaredLogger {
	InitLogClient()
	return sugar
}
