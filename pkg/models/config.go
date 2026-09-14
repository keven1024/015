package models

type Config struct {
	Share    ShareConfig    `mapstructure:"share"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Features FeaturesConfig `mapstructure:"features"`
	Site     SiteConfig     `mapstructure:"site"`
	About    AboutConfig    `mapstructure:"about"`
	SMTP     *SMTPConfig    `mapstructure:"smtp"`
}

type ShareConfig struct {
	DownloadSecret string `mapstructure:"download_secret" validate:"required"` // 必填，下载 JWT 签名密钥。
	DownloadWindow int    `mapstructure:"download_window"`                     // 单位：小时，默认值需在加载配置时设置为 12。
	PasswordSalt   string `mapstructure:"password_salt" validate:"required"`   // 必填，密码盐。
}

type UploadConfig struct {
	RemoveExpire int64  `mapstructure:"remove_expire"` // 单位：小时。
	Path         string `mapstructure:"path"`
	Maximum      string `mapstructure:"maximum"` // 保留人类可读容量，例如 100GiB。
}

type RedisConfig struct {
	URL string `mapstructure:"url"`
}

type FeaturesConfig struct {
	FileShare         *FeatureConfig `mapstructure:"file-share"`
	TextShare         *FeatureConfig `mapstructure:"text-share"`
	FileImageCompress *FeatureConfig `mapstructure:"file-image-compress"`
	FileImageConvert  *FeatureConfig `mapstructure:"file-image-convert"`
}

type FeatureConfig struct {
	Enabled *bool `mapstructure:"enabled"`
}

type SiteConfig struct {
	URL      string            `mapstructure:"url"`
	Title    map[string]string `mapstructure:"title"`
	Desc     map[string]string `mapstructure:"desc"`
	Icon     string            `mapstructure:"icon"`
	BgURL    string            `mapstructure:"bg_url"`
	EnableBg *bool             `mapstructure:"enable_bg"`
}

type AboutConfig struct {
	BgURL   string            `mapstructure:"bg_url" validate:"omitempty,http_url"`
	Content map[string]string `mapstructure:"content"` // 语言代码对应 Markdown 内容。
	Email   string            `mapstructure:"email" validate:"omitempty,email"`
	Name    string            `mapstructure:"name"`
	URL     string            `mapstructure:"url" validate:"omitempty,http_url"`
	Avatar  string            `mapstructure:"avatar"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Protocol string `mapstructure:"protocol"` // ssl 或 tls。
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
