package utils

import (
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	v       *viper.Viper
	envOnce sync.Once
)

func createViperInstance(props EnvOption) *viper.Viper {
	instance := viper.NewWithOptions(viper.ExperimentalBindStruct())
	instance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	instance.AutomaticEnv()
	for _, viperConfigType := range props.ConfigType {
		instance.SetConfigType(viperConfigType)
	}
	for _, name := range props.ConfigName {
		instance.SetConfigName(name)
	}
	for _, path := range props.ConfigPath {
		instance.AddConfigPath(path)
	}
	instance.OnConfigChange(func(_ fsnotify.Event) {
		LoadConfig()
	})
	instance.WatchConfig()
	if err := instance.ReadInConfig(); err != nil {
		panic(err)
	}
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

func SetEnv(key string, value string) {
	GetViperClient().Set(key, value)
}
