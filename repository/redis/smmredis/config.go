package smmredis

import "time"

type Config struct {
	SMMTTL time.Duration `koanf:"smm_ttl"`
}
