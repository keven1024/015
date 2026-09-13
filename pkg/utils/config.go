package utils

import (
	"pkg/models"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	config *models.Config
)

func loadConfig(v *viper.Viper) *models.Config {
	cfg := models.Config{
		Share:  models.ShareConfig{DownloadWindow: 12},
		Upload: models.UploadConfig{RemoveExpire: 2},
		Site: models.SiteConfig{
			Icon:     "/logo.png",
			BgURL:    "https://img.fudaoyuan.icu/api/1/random/?scale_min=1.5&webp=true&md=false&format=302",
			EnableBg: true,
		},
	}

	if err := v.Unmarshal(&cfg); err != nil {
		zap.L().Error("parse config error", zap.Error(err))
		return nil
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
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

func LoadConfig() {
	v := GetViperClient()
	config = loadConfig(v)
}
