package postgresproduct

import (
	"context"
	"fmt"
	"telegram-service-platform/entity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/repository/postgres"

	"github.com/shopspring/decimal"
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

		plan, sErr := scanStars(rows)

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

func scanStars(rows postgres.Scanner) (entity.StarPackage, error) {
	var plan entity.StarPackage
	var amountStr string
	sErr := rows.Scan(
		&plan.ID,
		&amountStr,
		&plan.Active,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)
	if sErr != nil {
		return entity.StarPackage{}, sErr
	}
	amount, err := decimal.NewFromString(amountStr)
	plan.Amount = entity.Amount(amount)
	if err != nil {
		return plan, fmt.Errorf("failed to parse amount to decimal: %w", err)
	}
	return plan, sErr
}
