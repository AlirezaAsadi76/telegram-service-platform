package orderfulfillmenttesting

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type fakeQueue struct {
	messages [][]string
	popIndex int

	queueKey string

	pushed []any

	events []string
}

func (f *fakeQueue) BRPop(
	_ context.Context,
	_ time.Duration,
	key string,
) ([]string, error) {
	f.events = append(f.events, "brpop")

	f.queueKey = key

	if f.popIndex >= len(f.messages) {
		return nil, redis.Nil
	}

	message := f.messages[f.popIndex]
	f.popIndex++

	return message, nil
}

func (f *fakeQueue) LPush(
	_ context.Context,
	key string,
	value any,
) error {
	f.events = append(f.events, "lpush")

	f.queueKey = key
	f.pushed = append(f.pushed, value)

	return nil
}
