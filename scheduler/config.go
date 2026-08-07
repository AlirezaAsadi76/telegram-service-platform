package scheduler

import "time"

type Config struct {
	currencyRefreshInterval time.Duration `koanf:"currency_refresh_interval"`
	starsRefreshInterval    time.Duration `koanf:"stars_refresh_interval"`
	premiumRefreshInterval  time.Duration `koanf:"premium_refresh_interval"`
}
