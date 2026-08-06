package fzrcards

import (
	"net/http"
)

type Config struct {
	BaseURL string `koanf:"base_url"`
	ApiKey  string `koanf:"api_key"`
}

type FzrClient struct {
	client *http.Client
	config Config
}

func New(cfg Config, client *http.Client) FzrClient {
	return FzrClient{
		client: client,
		config: cfg,
	}
}

func (a FzrClient) Connection() *http.Client {
	return a.client
}

func (a FzrClient) BaseURL() string {
	return a.config.BaseURL
}
