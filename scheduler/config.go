package scheduler

import "time"

type Config struct {
	CurrencyRefreshInterval time.Duration `koanf:"currency_refresh_interval"`
	StarsRefreshInterval    time.Duration `koanf:"stars_refresh_interval"`
	PremiumRefreshInterval  time.Duration `koanf:"premium_refresh_interval"`
}
