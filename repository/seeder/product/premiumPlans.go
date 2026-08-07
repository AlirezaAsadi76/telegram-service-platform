package product

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedPremiumPlans(
	ctx context.Context,
	db *pgxpool.Pool,
) error {

	query := `
	INSERT INTO premium_plans
		(duration)
	VALUES
		(3),
		(6),
		(12)
	ON CONFLICT DO NOTHING;
	`

	_, err := db.Exec(
		ctx,
		query,
	)

	return err
}
