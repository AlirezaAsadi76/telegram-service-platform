package orderfulfillmentrecoverytesting

import "context"

type fakeQueue struct {
	events []string

	queueKey string
	values   []any

	err error

	failOrderIDs map[uint64]error
}

func (f *fakeQueue) LPush(_ context.Context, queueKey string, value any) error {
	f.events = append(f.events, "enqueue")

	f.queueKey = queueKey
	f.values = append(f.values, value)

	orderID, ok := value.(uint64)
	if ok {
		if err, exists := f.failOrderIDs[orderID]; exists {
			return err
		}
	}

	return f.err
}
