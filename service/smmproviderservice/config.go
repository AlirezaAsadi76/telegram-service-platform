package smmproviderservice

import "time"

type Config struct {
	FailureThreshold int           `koanf:"failure-threshold"`
	SuccessThreshold int           `koanf:"success-threshold"`
	circuitTimeout   time.Duration `koanf:"circuit-timeout"`
}
