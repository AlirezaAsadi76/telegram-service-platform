package postgresproduct

import (
	"context"
	"telegram-service-platform/repository/postgres"

	"telegram-service-platform/entity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (db *DB) GetPremiumPlans(ctx context.Context) ([]entity.PremiumPlan, error) {

	const Op = "postgresproduct.GetPremiumPlans"

	query := `
		SELECT
			id,
			duration,
			active,
			created_at,
			updated_at
		FROM premium_plans
		WHERE active = TRUE
		ORDER BY duration ASC
	`

	rows, err := db.Pool.Query(
		ctx,
		query,
	)

	if err != nil {

		return nil,
			richerror.New(
				Op,
				err,
			).WithKind(richerror.KindQueryFailure).
				WithMessage(msgerror.QueryFailed)
	}

	defer rows.Close()

	plans := make([]entity.PremiumPlan, 0)

	for rows.Next() {

		plan, sErr := scanPremiumPlans(rows)

		if sErr != nil {
			return nil,
				richerror.New(
					Op,
					sErr,
				).WithKind(richerror.KindScanFailure).WithMessage(msgerror.QueryScanFailed)
		}
		plans = append(plans, plan)
	}

	if rErr := rows.Err(); rErr != nil {
		return nil, richerror.New(Op, rErr).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	return plans, nil
}

func scanPremiumPlans(rows postgres.Scanner) (entity.PremiumPlan, error) {
	var plan entity.PremiumPlan

	err := rows.Scan(
		&plan.ID,
		&plan.Duration,
		&plan.Active,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)
	return plan, err
}
