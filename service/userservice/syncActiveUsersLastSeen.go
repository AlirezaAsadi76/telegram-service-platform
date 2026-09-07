package userservice

import (
	"context"
	"telegram-service-platform/params/userparams"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) SyncActiveUsersLastSeen(ctx context.Context, req userparams.SyncLastSeenRequest) (userparams.SyncLastSeenResponse, error) {
	const op = "userservice.SyncActiveUsersLastSeen"

	activeIDs, err := s.activityTracker.GetActiveUsers(ctx)
	if err != nil {
		return userparams.SyncLastSeenResponse{}, richerror.New(op, err).
			WithKind(richerror.KindUnexpected)
	}

	if len(activeIDs) == 0 {
		return userparams.SyncLastSeenResponse{Synced: 0}, nil
	}

	rowsAffected, uErr := s.repository.UpdateLastSeenBulk(ctx, activeIDs)
	if uErr != nil {
		return userparams.SyncLastSeenResponse{}, richerror.New(op, uErr).
			WithKind(richerror.KindQueryFailure)
	}

	return userparams.SyncLastSeenResponse{Synced: rowsAffected}, nil
}
