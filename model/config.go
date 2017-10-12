package model

import (
	"time"
	"errors"
)

//Config configuration of application
type Config struct {
	FB           FBConfig       `json:"fb"`
	Telegram     TelegramConfig `json:"telegram"`
	DebugEnabled bool           `json:"debug_enabled"`
	Filter       Filter         `json:"filter"`
}

//FBConfig config of fb application and subscription
type FBConfig struct {
	AppSecret   string `json:"app_secret"`
	AppId       string `json:"app_id"`
	ClientToken string `json:"client_token"`

	DebugEnabled bool   `json:"debug_enabled"`
	DebugMode    string `json:"debug_mode"`

	Delay     time.Duration `json:"delay"`
	FeedURL   string        `json:"feed_url"`
	FeedLimit uint32        `json:"feed_limit"`
}

//TelegramConfig configuration of telegram application
type TelegramConfig struct {
	Token        string `json:"token"`
	DebugEnabled bool   `json:"debug_enabled"`
}

type Filter struct {
	KeywordsLandlord string `json:"keywords_landlord"`
	KeywordsTenant   string `json:"keywords_tenant"`
	KeywordsSpam     string `json:"keywords_spam"`
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