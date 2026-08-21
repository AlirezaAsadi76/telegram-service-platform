package userservice

import (
	"context"
	"telegram-service-platform/params"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) TrackActivity(ctx context.Context, req params.TrackUserActivityRequest) (params.TrackUserActivityResponse, error) {
	const op = "userservice.TrackActivity"

	if req.TelegramID <= 0 {
		return params.TrackUserActivityResponse{}, richerror.New(op, nil).
			WithKind(richerror.KindValidation).
			WithMessage(msgerror.InvalidInput)
	}

	if err := s.activityTracker.TrackActivity(ctx, req.TelegramID.Int64()); err != nil {
		return params.TrackUserActivityResponse{}, richerror.New(op, err).
			WithKind(richerror.KindUnexpected)
	}

	return params.TrackUserActivityResponse{Tracked: true}, nil
}
