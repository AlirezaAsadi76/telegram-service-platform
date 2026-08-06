package exchangerate

import "time"

type Config struct {
	TonUsdURL string        `koanf:"ton_usd_url"`
	UsdIrURL  string        `koanf:"usd_irr_url"`
	Timeout   time.Duration `koanf:"timeout"`
}
