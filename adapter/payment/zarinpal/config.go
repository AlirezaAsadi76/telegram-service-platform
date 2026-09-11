package zarinpal

import "time"

type Config struct {
	MerchantID  string        `koanf:"merchant_id"`
	BaseURL     string        `koanf:"base_url"`
	Timeout     time.Duration `koanf:"timeout"`
	StartPayURL string        `koanf:"start_pay_url"`
}
