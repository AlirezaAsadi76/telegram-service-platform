package justanotherpanel

import "time"

type Config struct {
	BaseURL string        `koanf:"base_url"`
	APIKey  string        `koanf:"api_key"`
	Timeout time.Duration `koanf:"timeout"`
}
