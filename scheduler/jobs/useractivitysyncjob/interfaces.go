package useractivitysyncjob

import (
	"context"
	"telegram-service-platform/params"
)

type UserService interface {
	SyncActiveUsersLastSeen(ctx context.Context, req params.SyncLastSeenRequest) (params.SyncLastSeenResponse, error)
}
