package orderfulfillmentservice

import "context"

type QueueRepository interface {
	LPush(ctx context.Context, queueKey string, data any) error
}
