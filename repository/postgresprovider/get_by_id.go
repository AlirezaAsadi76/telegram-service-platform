package postgresprovider

import (
	"context"
	"errors"

	"telegram-service-platform/entity/providerentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"github.com/jackc/pgx/v5"
)

func (d *DB) GetByID(ctx context.Context, providerID uint64) (*providerentity.Provider, error) {
	const Op = "postgresprovider.GetByID"

	query := `
		SELECT
			id,
			name,
			type,
			base_url,
			api_key,
			config,
			priority,
			is_active
		FROM providers
		WHERE id = $1
	`

	provider, err := scanProvider(
		d.Pool.Connection().QueryRow(ctx, query, providerID),
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, richerror.New(Op, err).
				WithKind(richerror.KindNotFound).
				WithMessage(msgerror.NoAvailableAdapter)
		}

		return nil, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	return &provider, nil
}
