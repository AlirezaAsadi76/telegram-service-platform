package priceservice

import "time"

type Config struct {
	StarsPriceTTL   time.Duration `koanf:"stars_price_ttl"`
	PremiumPriceTTL time.Duration `koanf:"premium_price_ttl"`
	CurrencyTTL     time.Duration `koanf:"currency_ttl"`
}
