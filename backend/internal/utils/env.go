package utils

import (
	"github.com/spf13/viper"
)

var v = viper.New()

func init() {
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
}

func GetEnv(key string) string {
	return v.GetString(key)
}
