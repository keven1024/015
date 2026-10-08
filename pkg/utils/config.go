package utils

import (
	"pkg/models"

	"dario.cat/mergo"
	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	config *models.Config
)

func defaultConfig() models.Config {
	return models.Config{
		Share:  models.ShareConfig{DownloadWindow: 12},
		Upload: models.UploadConfig{RemoveExpire: 2},
		Site: models.SiteConfig{
			Icon:     "/logo.webp",
			BgURL:    "https://img.fudaoyuan.icu/api/1/random/?scale_min=1.5&webp=true&md=false&format=302",
			EnableBg: lo.ToPtr(true),
		},
		Features: models.FeaturesConfig{
			models.FeatureConfigKeyFileShare: map[string]any{
				"enabled": true,
			},
			models.FeatureConfigKeyTextShare: map[string]any{
				"enabled": true,
			},
		},
	}
}

func loadConfig(v *viper.Viper) *models.Config {
	cfg := defaultConfig()
	if err := v.Unmarshal(&cfg); err != nil {
		zap.L().Error("parse config error", zap.Error(err))
		return nil
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("storage_size", func(fl validator.FieldLevel) bool {
		size, err := GetFileSize(fl.Field().String())
		return err == nil && size > 0
	})
	validate.RegisterValidation("required_en_key_content", func(fl validator.FieldLevel) bool {
		return fl.Field().Interface().(models.MultiContent)["en"] != ""
	})
	if err := validate.Struct(cfg); err != nil {
		zap.L().Error("validate config error", zap.Error(err))
		return nil
	}
	zap.L().Info("load config ok")
	return &cfg
}

func GetConfig() *models.Config {
	return config
}

// SetConfig 用默认配置补齐零值字段并替换当前配置，不执行校验。
func SetConfig(cfg models.Config) {
	if err := mergo.Merge(&cfg, defaultConfig()); err != nil {
		panic(err)
	}
	config = &cfg
}

func LoadConfig() {
	v := GetViperClient()
	config = loadConfig(v)
}
