package statussynctesting

import (
	"context"

	"telegram-service-platform/entity/notificationentity"
)

type fakeNotificationRepository struct {
	created []*notificationentity.Notification

	err error
}

func (f *fakeNotificationRepository) Create(
	_ context.Context,
	notification *notificationentity.Notification,
) error {
	if f.err != nil {
		return f.err
	}

	f.created = append(
		f.created,
		notification,
	)

	return nil
}

func (f *fakeNotificationRepository) GetPending(
	_ context.Context,
	_ int,
) ([]notificationentity.Notification, error) {
	return nil, nil
}

func (f *fakeNotificationRepository) UpdateStatus(
	_ context.Context,
	_ uint64,
	_ notificationentity.NotificationStatus,
) error {
	return nil
}

func (f *fakeNotificationRepository) UpdateRetryCount(
	_ context.Context,
	_ uint64,
	_ int,
) error {
	return nil
}

type fakeNotificationRedis struct {
	enqueued []any
}

func (f *fakeNotificationRedis) LPush(
	_ context.Context,
	_ string,
	data any,
) error {
	f.enqueued = append(
		f.enqueued,
		data,
	)

	return nil
}
