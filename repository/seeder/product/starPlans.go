package product

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedStarPlans(
	ctx context.Context,
	db *pgxpool.Pool,
) error {

	query := `
	INSERT INTO star_plans
		(amount)
	VALUES
		(50),
		(100),
		(250),
		(500),
		(1000)
	ON CONFLICT DO NOTHING;
	`

	_, err := db.Exec(
		ctx,
		query,
	)

	return err
}
