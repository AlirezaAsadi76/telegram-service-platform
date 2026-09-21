package orderfulfillmenttesting

import (
	"context"
	"telegram-service-platform/entity/notificationentity"
)

type fakeNotificationStore struct {
	created  []*notificationentity.Notification
	enqueued []any

	events []string
}

func (f *fakeNotificationStore) Create(
	_ context.Context,
	notification *notificationentity.Notification,
) error {
	f.events = append(
		f.events,
		"notification_create",
	)

	f.created = append(
		f.created,
		notification,
	)

	return nil
}

func (f *fakeNotificationStore) GetPending(
	_ context.Context,
	_ int,
) ([]notificationentity.Notification, error) {
	return nil, nil
}

func (f *fakeNotificationStore) UpdateStatus(
	_ context.Context,
	_ uint64,
	_ notificationentity.NotificationStatus,
) error {
	return nil
}

func (f *fakeNotificationStore) UpdateRetryCount(
	_ context.Context,
	_ uint64,
	_ int,
) error {
	return nil
}

func (f *fakeNotificationStore) LPush(
	_ context.Context,
	_ string,
	value any,
) error {
	f.events = append(
		f.events,
		"notification_enqueue",
	)

	f.enqueued = append(
		f.enqueued,
		value,
	)

	return nil
}
