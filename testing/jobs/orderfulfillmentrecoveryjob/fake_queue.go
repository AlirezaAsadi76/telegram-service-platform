package orderfulfillmentrecoverytesting

import "context"

type fakeQueue struct {
	events []string

	queueKey string
	values   []any

	err error
}

func (f *fakeQueue) LPush(
	_ context.Context,
	queueKey string,
	value any,
) error {
	f.events = append(f.events, "enqueue")
	f.queueKey = queueKey
	f.values = append(f.values, value)

	return f.err
}
