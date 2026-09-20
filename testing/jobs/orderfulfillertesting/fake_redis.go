package orderfulfillertesting

import (
	"context"
	"time"
)

type fakeRedis struct {
	result []string
	err    error

	events []string
}

func (f *fakeRedis) BRPop(_ context.Context, _ time.Duration, _ string) ([]string, error) {
	f.events = append(f.events, "queue_pop")

	return f.result, f.err
}
