package fzrcards

import (
	"net/http"
)

type FzrClient struct {
	client *http.Client
	config Config
}

func New(cfg Config) *FzrClient {
	client := &http.Client{
		Timeout: cfg.Timeout,
	}
	return &FzrClient{
		client: client,
		config: cfg,
	}
}

func (a FzrClient) Connection() *http.Client {
	return a.client
}
