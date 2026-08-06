package exchangerate

type Config struct {
	TonUsdURL string `koanf:"ton_usd_url"`
	UsdIrURL  string `koanf:"usd_irr_url"`
}
