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
	DownloadWindow int    `mapstructure:"download_window" validate:"gt=0"`     // 单位：小时，默认值需在加载配置时设置为 12。
	PasswordSalt   string `mapstructure:"password_salt" validate:"required"`   // 必填，密码盐。
}

type UploadConfig struct {
	RemoveExpire int64  `mapstructure:"remove_expire" validate:"gt=0"` // 单位：小时。
	Path         string `mapstructure:"path" validate:"dirpath"`
	Maximum      string `mapstructure:"maximum" validate:"required,storage_size"` // 保留人类可读容量，例如 100GiB。
}

type RedisConfig struct {
	URL string `mapstructure:"url" validate:"required"`
}

type FeatureConfigKey string
type FeaturesConfig map[FeatureConfigKey]any

const (
	FeatureConfigKeyFileShare         FeatureConfigKey = "file-share"
	FeatureConfigKeyTextShare         FeatureConfigKey = "text-share"
	FeatureConfigKeyFileImageCompress FeatureConfigKey = "file-image-compress"
	FeatureConfigKeyFileImageConvert  FeatureConfigKey = "file-image-convert"
)

type FeatureConfigBaseValue struct {
	Enabled *bool `mapstructure:"enabled"`
}

type SiteConfig struct {
	URL      string       `mapstructure:"url" validate:"required,http_url"`
	Title    MultiContent `mapstructure:"title" validate:"required,required_en_key_content"`
	Desc     MultiContent `mapstructure:"desc" validate:"required,required_en_key_content"`
	Icon     string       `mapstructure:"icon"`
	BgURL    string       `mapstructure:"bg_url" validate:"omitempty,http_url"`
	EnableBg *bool        `mapstructure:"enable_bg"`
}

type AboutConfig struct {
	BgURL   string       `mapstructure:"bg_url" validate:"omitempty,http_url"`
	Content MultiContent `mapstructure:"content" validate:"omitempty,required_en_key_content"` // 语言代码对应 Markdown 内容。
	Email   string       `mapstructure:"email" validate:"omitempty,email"`
	Name    string       `mapstructure:"name"`
	URL     string       `mapstructure:"url" validate:"omitempty,http_url"`
	Avatar  string       `mapstructure:"avatar" validate:"omitempty,http_url"`
}

type MultiContent map[string]string

type SMTPConfig struct {
	Host     string `mapstructure:"host" validate:"omitempty,hostname_rfc1123|ip"`
	Port     uint   `mapstructure:"port" validate:"omitempty,port"`
	Protocol string `mapstructure:"protocol" validate:"omitempty,oneof=ssl tls"` // ssl 或 tls。
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
