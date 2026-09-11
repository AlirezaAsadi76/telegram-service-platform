package zarinpal

import "net/http"

type Adapter struct {
	config Config
	client *http.Client
}

func New(config Config, client *http.Client) *Adapter {
	if client == nil {
		client = &http.Client{
			Timeout: config.Timeout,
		}
	}

	return &Adapter{
		config: config,
		client: client,
	}
}
