package postgresprovider

import (
	"context"
	"telegram-service-platform/entity/providerentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) GetActiveByType(ctx context.Context, providerType providerentity.ProviderType) ([]*providerentity.Provider, error) {
	const Op = "postgresprovider.GetActiveByType"

	query := `
		SELECT id, name, type, base_url, api_key, config, priority, is_active, created_at
		FROM providers WHERE type = $1 AND is_active = true ORDER BY priority ASC
	`
	rows, qErr := d.Pool.Connection().Query(ctx, query, providerType)
	if qErr != nil {
		return nil, richerror.New(Op, qErr).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}
	defer rows.Close()

	providers := make([]*providerentity.Provider, 0)

	for rows.Next() {

		prov, sErr := scanProvider(rows)
		if sErr != nil {
			return nil, richerror.New(Op, sErr).WithKind(richerror.KindScanFailure).WithMessage(msgerror.QueryScanFailed)
		}

		providers = append(
			providers,
			&prov,
		)
	}

	return providers, rows.Err()
}
