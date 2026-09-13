package utils

import (
	"pkg/models"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetConfigMergesDefaultsAndReplacesPreviousValues(t *testing.T) {
	previous := config
	t.Cleanup(func() { config = previous })
	SetConfig(models.Config{
		Share: models.ShareConfig{PasswordSalt: "testsalt", DownloadWindow: 6},
		Site:  models.SiteConfig{Icon: "/custom.png"},
	})
	assert.Equal(t, "testsalt", GetConfig().Share.PasswordSalt)
	assert.Equal(t, 6, GetConfig().Share.DownloadWindow)
	assert.Equal(t, "/custom.png", GetConfig().Site.Icon)
	assert.EqualValues(t, 2, GetConfig().Upload.RemoveExpire)
	SetConfig(models.Config{Upload: models.UploadConfig{Path: "/tmp/uploads"}})
	assert.Equal(t, "/tmp/uploads", GetConfig().Upload.Path)
	assert.Empty(t, GetConfig().Share.PasswordSalt)
	assert.Equal(t, 12, GetConfig().Share.DownloadWindow)
	assert.Equal(t, "/logo.png", GetConfig().Site.Icon)
	assert.True(t, GetConfig().Site.EnableBg)
}

func TestLoadConfigDefaultsAndOverrides(t *testing.T) {
	t.Setenv("SHARE_DOWNLOAD_SECRET", "testsecret")
	t.Setenv("SHARE_PASSWORD_SALT", "testsalt")
	t.Setenv("UPLOAD_PATH", "/tmp/uploads")
	instance := viper.NewWithOptions(viper.ExperimentalBindStruct())
	instance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	instance.AutomaticEnv()
	instance.Set("site.enable_bg", false)
	instance.Set("share.download_window", 6)
	cfg := loadConfig(instance)
	require.NotNil(t, cfg)
	assert.Equal(t, "/tmp/uploads", cfg.Upload.Path)
	assert.EqualValues(t, 2, cfg.Upload.RemoveExpire)
	assert.Equal(t, 6, cfg.Share.DownloadWindow)
	assert.Equal(t, "/logo.png", cfg.Site.Icon)
	assert.False(t, cfg.Site.EnableBg)
	instance.Set("share.download_window", 8)
	assert.Equal(t, 8, loadConfig(instance).Share.DownloadWindow)
}
