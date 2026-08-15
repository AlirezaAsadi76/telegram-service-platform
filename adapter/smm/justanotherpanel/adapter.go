package justanotherpanel

import (
	"net/http"
)

type Adapter struct {
	client *http.Client
	config Config
}

func New(cfg Config) *Adapter {

	return &Adapter{
		client: &http.Client{Timeout: cfg.Timeout},
		config: cfg,
	}
}
