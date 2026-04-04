package utils

import (
	"io"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	v       *viper.Viper
	envOnce sync.Once
)

func createViperInstance(props EnvOption) *viper.Viper {
	instance := viper.New()
	for _, viperConfigType := range props.ConfigType {
		instance.SetConfigType(viperConfigType)
	}
	if props.ConfigData != nil {
		instance.ReadConfig(props.ConfigData)
		return instance
	}
	for _, name := range props.ConfigName {
		instance.SetConfigName(name)
	}
	instance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	for _, path := range props.ConfigPath {
		instance.AddConfigPath(path)
	}
	instance.AutomaticEnv()
	instance.WatchConfig()
	if err := instance.ReadInConfig(); err != nil {
		panic(err)
	}
	return instance
}

func InitTestViper(props EnvOption) *viper.Viper {
	instance := createViperInstance(props)
	v = instance
	envOnce.Do(func() {}) // 消费 once，防止 GetViperClient 覆盖已注入的实例
	return instance
}

func GetViperClient() *viper.Viper {
	envOnce.Do(func() {
		v = createViperInstance(getEnvOptions())
	})
	return v
}

type Option interface {
	applyTo(*EnvOption)
}

type EnvOption struct {
	DefaultValue string
	ConfigPath   []string
	ConfigName   []string
	ConfigType   []string
	ConfigData   io.Reader // 测试环境使用
}

type WithDefaultValue string

func (o WithDefaultValue) applyTo(props *EnvOption) {
	props.DefaultValue = string(o)
}

func getEnvOptions(options ...Option) EnvOption {
	props := EnvOption{
		DefaultValue: "",
		ConfigPath:   []string{".", "../"},
		ConfigName:   []string{"config"},
		ConfigType:   []string{"yaml"},
	}
	for _, option := range options {
		option.applyTo(&props)
	}
	return props
}

func GetEnv(key string, options ...Option) string {
	props := getEnvOptions(options...)
	value := GetViperClient().GetString(key)

	if value == "" && props.DefaultValue != "" {
		return props.DefaultValue
	}
	return value
}

func GetEnvWithDefault(key string, defaultValue string) string {
	return GetEnv(key, WithDefaultValue(defaultValue))
}

func GetEnvMap(key string) map[string]any {
	return GetViperClient().GetStringMap(key)
}

func SetEnv(key string, value string) {
	GetViperClient().Set(key, value)
}
