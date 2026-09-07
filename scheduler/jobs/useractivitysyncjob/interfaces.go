package useractivitysyncjob

import (
	"context"
	"telegram-service-platform/params/userparams"
)

type UserService interface {
	SyncActiveUsersLastSeen(ctx context.Context, req userparams.SyncLastSeenRequest) (userparams.SyncLastSeenResponse, error)
}
