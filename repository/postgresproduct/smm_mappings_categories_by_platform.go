package postgresproduct

import (
	"context"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (db *DB) SMMMappingGetDistinctCategoriesByPlatform(ctx context.Context, platform smmentity.PlatformType) ([]smmentity.Category, error) {
	const op = "postgresproduct.SMMMappingGetDistinctCategoriesByPlatform"

	query := `
		SELECT DISTINCT category 
		FROM smm_service_mappings 
		WHERE platform = $1 AND is_active = true 
		ORDER BY category
	`

	rows, qErr := db.Pool.Connection().Query(ctx, query, platform)
	if qErr != nil {
		return nil, richerror.New(op, qErr).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}
	defer rows.Close()

	var categories []smmentity.Category
	for rows.Next() {
		var c smmentity.Category
		if err := rows.Scan(&c.Name); err != nil {
			return nil, richerror.New(op, err).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryScanFailed)
		}
		categories = append(categories, c)
	}

	return categories, nil
}
