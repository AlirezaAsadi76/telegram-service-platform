package postgresproduct

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (db *DB) GetStarPlans(ctx context.Context) ([]entity.StarPackage, error) {
	const Op = "postgresproduct.GetStarPlans"

	query := `
		SELECT
			id,
			amount,
			active,
			created_at,
			updated_at
		FROM star_plans
		WHERE active = true
		ORDER BY amount ASC
	`

	rows, qErr := db.Pool.Connection().Query(
		ctx,
		query,
	)

	if qErr != nil {
		return nil,
			richerror.New(Op, qErr).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}

	defer rows.Close()

	var plans []entity.StarPackage

	for rows.Next() {

		var plan entity.StarPackage

		sErr := rows.Scan(
			&plan.ID,
			&plan.Amount,
			&plan.Active,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)

		if sErr != nil {

			return nil, richerror.New(Op, sErr).
				WithKind(richerror.KindUnexpected).
				WithMessage(msgerror.Unexpected)
		}

		plans = append(plans, plan)
	}

	if err := rows.Err(); err != nil {

		return nil, richerror.New(Op, err).WithKind(richerror.KindUnexpected).WithMessage(msgerror.Unexpected)
	}

	return plans, nil
}
