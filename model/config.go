package model

//Config configuration of application
type Config struct {
	FB       FBConfig       `json:"fb"`
	Telegram TelegramConfig `json:"telegram"`
}

//FBConfig config of fb application and subscription
type FBConfig struct {
	AppSecret   string `json:"app_secret"`
	AppId       string `json:"app_id"`
	ClientToken string `json:"client_token"`

	FeedURL string `json:"feed_url"`
}

//TelegramConfig configuration of telegram application
type TelegramConfig struct {
}
