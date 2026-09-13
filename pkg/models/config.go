package models

type Config struct {
	Share    ShareConfig    `yaml:"share" mapstructure:"share"`
	Upload   UploadConfig   `yaml:"upload" mapstructure:"upload"`
	Redis    RedisConfig    `yaml:"redis" mapstructure:"redis"`
	Features FeaturesConfig `yaml:"features" mapstructure:"features"`
	Site     SiteConfig     `yaml:"site" mapstructure:"site"`
	About    AboutConfig    `yaml:"about" mapstructure:"about"`
	SMTP     SMTPConfig     `yaml:"smtp" mapstructure:"smtp"`
}

type ShareConfig struct {
	DownloadSecret string `yaml:"download_secret" mapstructure:"download_secret" validate:"required"` // 必填，下载 JWT 签名密钥。
	DownloadWindow int    `yaml:"download_window" mapstructure:"download_window"`                     // 单位：小时，默认值需在加载配置时设置为 12。
	PasswordSalt   string `yaml:"password_salt" mapstructure:"password_salt" validate:"required"`     // 必填，密码盐。
}

type UploadConfig struct {
	RemoveExpire int64  `yaml:"remove_expire" mapstructure:"remove_expire"` // 单位：小时。
	Path         string `yaml:"path" mapstructure:"path"`
	Maximum      string `yaml:"maximum" mapstructure:"maximum"` // 保留人类可读容量，例如 100GiB。
}

type RedisConfig struct {
	URL string `yaml:"url" mapstructure:"url"`
}

type FeaturesConfig struct {
	FileShare         FeatureConfig `yaml:"file-share" mapstructure:"file-share"`
	TextShare         FeatureConfig `yaml:"text-share" mapstructure:"text-share"`
	FileImageCompress FeatureConfig `yaml:"file-image-compress" mapstructure:"file-image-compress"`
	FileImageConvert  FeatureConfig `yaml:"file-image-convert" mapstructure:"file-image-convert"`
}

type FeatureConfig struct {
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
}

type SiteConfig struct {
	URL      string            `yaml:"url" mapstructure:"url"`
	Title    map[string]string `yaml:"title" mapstructure:"title"`
	Desc     map[string]string `yaml:"desc" mapstructure:"desc"`
	Icon     string            `yaml:"icon" mapstructure:"icon"`
	BgURL    string            `yaml:"bg_url" mapstructure:"bg_url"`
	EnableBg bool              `yaml:"enable_bg" mapstructure:"enable_bg"`
}

type AboutConfig struct {
	BgURL   string            `yaml:"bg_url" mapstructure:"bg_url"`
	Content map[string]string `yaml:"content" mapstructure:"content"` // 语言代码对应 Markdown 内容。
	Email   string            `yaml:"email" mapstructure:"email"`
	Name    string            `yaml:"name" mapstructure:"name"`
	URL     string            `yaml:"url" mapstructure:"url"`
	Avatar  string            `yaml:"avatar" mapstructure:"avatar"`
}

type SMTPConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	Protocol string `yaml:"protocol" mapstructure:"protocol"` // ssl 或 tls。
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
}
