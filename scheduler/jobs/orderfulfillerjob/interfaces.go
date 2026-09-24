package orderfulfillerjob

import (
	"context"
	"telegram-service-platform/entity/smmentity"
	"time"
)

type RedisRepository interface {
	BRPop(ctx context.Context, ttl time.Duration, key string) ([]string, error)
}

type ProductService interface {
	GetSMMServiceByMappingID(ctx context.Context, mappingID int64) (*smmentity.SMM, error)
}
