package postgresproduct

import (
	"context"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (db *DB) SMMMappingGetDistinctPlatforms(ctx context.Context) ([]smmentity.Platform, error) {
	const op = "postgresproduct.SMMMappingGetDistinctPlatforms"

	query := `
		SELECT DISTINCT platform 
		FROM smm_service_mappings 
		WHERE is_active = true 
		ORDER BY platform
	`

	rows, qErr := db.Pool.Connection().Query(ctx, query)
	if qErr != nil {
		return nil, richerror.New(op, qErr).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}
	defer rows.Close()

	var platforms []smmentity.Platform
	for rows.Next() {
		var p smmentity.Platform
		if err := rows.Scan(&p.Name); err != nil {
			return nil, richerror.New(op, err).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
		}
		platforms = append(platforms, p)
	}

	return platforms, nil
}
