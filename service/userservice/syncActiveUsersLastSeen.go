package userservice

import (
	"context"
	"telegram-service-platform/params"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) SyncActiveUsersLastSeen(ctx context.Context, req params.SyncLastSeenRequest) (params.SyncLastSeenResponse, error) {
	const op = "userservice.SyncActiveUsersLastSeen"

	activeIDs, err := s.activityTracker.GetActiveUsers(ctx)
	if err != nil {
		return params.SyncLastSeenResponse{}, richerror.New(op, err).
			WithKind(richerror.KindUnexpected)
	}

	if len(activeIDs) == 0 {
		return params.SyncLastSeenResponse{Synced: 0}, nil
	}

	rowsAffected, uErr := s.repository.UpdateLastSeenBulk(ctx, activeIDs)
	if uErr != nil {
		return params.SyncLastSeenResponse{}, richerror.New(op, uErr).
			WithKind(richerror.KindQueryFailure)
	}

	return params.SyncLastSeenResponse{Synced: rowsAffected}, nil
}
