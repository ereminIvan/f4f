package model

import (
	"time"
	"promo_tools_api/common/errors"
)

//Config configuration of application
type Config struct {
	FB             FBConfig       `json:"fb"`
	Telegram       TelegramConfig `json:"telegram"`
	DebugEnabled   bool           `json:"debug_enabled"`
	KeywordsFilter string         `json:"keywords_filter"`
	KeywordsSearch string         `json:"keywords_search"`
}

//FBConfig config of fb application and subscription
type FBConfig struct {
	AppSecret   string `json:"app_secret"`
	AppId       string `json:"app_id"`
	ClientToken string `json:"client_token"`

	DebugEnabled bool   `json:"debug_enabled"`
	DebugMode    string `json:"debug_mode"`

	FeedRequestFrequency time.Duration `json:"feed_request_frequency"`
	FeedURL              string        `json:"feed_url"`
}

//TelegramConfig configuration of telegram application
type TelegramConfig struct {
	Token        string `json:"token"`
	DebugEnabled bool   `json:"debug_enabled"`
}

//Valida config parameters
func (c *Config) Validate() (bool, error) {
	if c.FB.AppSecret == "" || c.FB.AppId == "" {
		return false, errors.New("Error: Facebook `app_secret`, `app_id`, `feed_url` could not be empty")
	}
	if c.Telegram.Token == "" {
		return false, errors.New("Error: Telegram `token` could not be empty")
	}
	return true, nil
}