package postgresproduct

import (
	"context"
	"telegram-service-platform/repository/postgres"

	"telegram-service-platform/entity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (db *DB) GetAdsPlans(ctx context.Context) ([]entity.AdsPlan, error) {
	const Op = "postgresproduct.GetAdsPlans"

	query := `
		SELECT
			id,
			views,
			cpm,
			daily_view_limit,
			active,
			created_at,
			updated_at
		FROM ads_plans
		WHERE active = TRUE
		ORDER BY views ASC
	`

	rows, err := db.Pool.Connection().Query(ctx, query)
	if err != nil {
		return nil, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	defer rows.Close()

	plans := make([]entity.AdsPlan, 0)

	for rows.Next() {
		plan, sErr := scanAdsPlan(rows)
		if sErr != nil {
			return nil, richerror.New(Op, err).
				WithKind(richerror.KindScanFailure).
				WithMessage(msgerror.QueryScanFailed)
		}

		plans = append(plans, plan)
	}

	if err := rows.Err(); err != nil {
		return nil, richerror.New(Op, err).
			WithKind(richerror.KindUnexpected).
			WithMessage(msgerror.Unexpected)
	}

	return plans, nil
}

func scanAdsPlan(rows postgres.Scanner) (entity.AdsPlan, error) {
	var plan entity.AdsPlan
	err := rows.Scan(
		&plan.ID,
		&plan.Views,
		&plan.CPM,
		&plan.DailyViewLimit,
		&plan.Active,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)

	return plan, err
}
