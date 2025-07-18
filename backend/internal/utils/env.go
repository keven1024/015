package utils

import (
	"github.com/spf13/viper"
)

var v *viper.Viper

func init() {
	InitEnv()
}

func InitEnv() {
	if v != nil {
		return
	}
	v = viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("../")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// 只有当错误不是"配置文件未找到"时才 panic
			panic(err)
		}
	}
}

func GetEnv(key string) string {
	InitEnv()
	return v.GetString(key)
}

func GetEnvWithDefault(key string, defaultValue string) string {
	value := v.GetString(key)
	if value == "" {
		return defaultValue
	}
	return value
}
