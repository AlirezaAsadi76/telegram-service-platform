package fzrcards

import "time"

type Config struct {
	BaseURL string        `koanf:"base_url"`
	ApiKey  string        `koanf:"api_key"`
	Timeout time.Duration `koanf:"timeout"`
}

func (a FzrClient) BaseURL() string {
	return a.config.BaseURL
}

func (a FzrClient) APIKey() string {
	return a.config.ApiKey
}
