package postgresorder

import (
	"context"
	"encoding/json"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) Create(ctx context.Context, order *orderentity.Order) error {
	const Op = "postgresorder.CreateOrder"
	metadata, err := json.Marshal(order.Metadata)
	if err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindSerializationFailure).WithMessage(msgerror.MarshalFailed)
	}

	query := `
		INSERT INTO orders (
			user_id,
			product_type,
			product_id,
			quantity,
			target_link,
			amount,
			currency,
			status,
			external_order_id,
			provider_id,
			metadata
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11
		)
		RETURNING id, created_at, updated_at
	`

	err = d.Pool.Connection().QueryRow(
		ctx,
		query,
		order.UserID,
		order.ProductType,
		order.ProductID,
		order.Quantity,
		order.TargetLink,
		order.Amount,
		order.Currency,
		order.Status,
		order.ExternalOrderID,
		order.ProviderID,
		metadata,
	).Scan(
		&order.ID,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}

	return nil
}
