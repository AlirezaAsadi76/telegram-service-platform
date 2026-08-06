package fzrcards

type Config struct {
	BaseURL string `koanf:"base_url"`
	ApiKey  string `koanf:"api_key"`
}

func (a FzrClient) BaseURL() string {
	return a.config.BaseURL
}

func (a FzrClient) APIKey() string {
	return a.config.ApiKey
}
