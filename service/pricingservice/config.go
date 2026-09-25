package pricingservice

type Config struct {
	SMMMinimumUnitPriceToman string `koanf:"smm_minimum_unit_price_toman"`
	SMMMinimumMultiplier     string `koanf:"smm_minimum_multiplier"`
	SMMMaximumMultiplier     string `koanf:"smm_maximum_multiplier"`
	SMMMultiplierScaleToman  string `koanf:"smm_multiplier_scale_toman"`
}
