package share

const modelName = "015:shareInfoMap"

type RedisShareInfo struct {
	// Id          string    `json:"id"`
	CreatedAt          int64           `json:"created_at"`
	UpdatedAt          int64           `json:"updated_at"`
	Owner              string          `json:"owner"`
	Type               ShareType       `json:"type"`
	Text               string          `json:"text"`
	Files              []ShareFileData `json:"files"`
	ExpireAt           int64           `json:"expire_time"`
	ViewNum            int64           `json:"download_nums"`
	Password           string          `json:"password"`
	NotifyEmails       []string        `json:"notify_emails"`
	NotifyWebhooks     []NotifyWebhook `json:"notify_webhooks"`
	Locale             string          `json:"locale"`
	PickupCode         string          `json:"pickup_code"`
	PickupCodeExpireAt int64           `json:"pickup_code_expire_at"`
}

type ShareFileData struct {
	Id       string `json:"id"`
	FileName string `json:"file_name"`
}

type NotifyWebhook struct {
	URL      string            `json:"url"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers"`
	BodyType string            `json:"bodyType"`
	Body     string            `json:"body"`
}

type ShareType string

const (
	ShareTypeFile ShareType = "file"
	ShareTypeText ShareType = "text"
)
