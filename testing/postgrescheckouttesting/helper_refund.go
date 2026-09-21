package postgrescheckouttesting

import (
	"context"
	"fmt"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/entity/providerentity"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type refund struct {
	UserID     uint64
	WalletID   uint64
	OrderID    uint64
	ProviderID uint64
}

func createRefundFixture(t *testing.T, pool *pgxpool.Pool, orderStatus orderentity.OrderStatus, amount decimal.Decimal) refund {
	t.Helper()

	ctx := context.Background()

	unique := time.Now().UnixNano()

	username := fmt.Sprintf(
		"refund-test-%d",
		unique,
	)

	telegramID := unique

	var userID uint64

	err := pool.QueryRow(
		ctx,
		`
		INSERT INTO users (
			telegram_id,
			first_name,
			last_name,
			username,
			role
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`,
		telegramID,
		"Refund",
		"Test",
		username,
		entity.UserRole,
	).Scan(&userID)

	if err != nil {
		t.Fatalf(
			"create test user: %v",
			err,
		)
	}

	var providerID uint64

	err = pool.QueryRow(
		ctx,
		`
		INSERT INTO providers (
			name,
			type,
			base_url,
			api_key,
			config,
			priority,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
		`,
		fmt.Sprintf("refund-test-provider-%d", unique),
		providerentity.ProviderTypeSMM,
		"https://provider.test",
		"test-api-key",
		[]byte(`{}`),
		1,
		true,
	).Scan(&providerID)

	if err != nil {
		t.Fatalf(
			"create test provider: %v",
			err,
		)
	}

	var walletID uint64

	err = pool.QueryRow(
		ctx,
		`
		INSERT INTO wallets (
			user_id,
			balance,
			currency,
			version
		)
		VALUES ($1, $2, $3, 1)
		RETURNING id
		`,
		userID,
		entity.Amount(decimal.NewFromInt(100)),
		entity.Currency("TON"),
	).Scan(&walletID)

	if err != nil {
		t.Fatalf(
			"create test wallet: %v",
			err,
		)
	}

	var orderID uint64

	err = pool.QueryRow(
		ctx,
		`
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
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11
		)
		RETURNING id
		`,
		userID,
		productentity.ProductTypeSMM,
		uint64(1),
		int64(100),
		"https://example.com/test",
		entity.Amount(amount),
		entity.Currency("TON"),
		orderStatus,
		fmt.Sprintf("EXT-REFUND-%d", unique),
		providerID,
		[]byte(`{}`),
	).Scan(&orderID)

	if err != nil {
		t.Fatalf(
			"create test order: %v",
			err,
		)
	}

	fixture := refund{
		UserID:     userID,
		WalletID:   walletID,
		OrderID:    orderID,
		ProviderID: providerID,
	}

	t.Cleanup(func() {
		cleanupRefundFixture(pool, fixture)
	})

	return fixture
}

func cleanupRefundFixture(pool *pgxpool.Pool, fixture refund) {
	ctx := context.Background()

	_, _ = pool.Exec(
		ctx,
		`DELETE FROM wallet_transactions WHERE wallet_id = $1`,
		fixture.WalletID,
	)

	_, _ = pool.Exec(
		ctx,
		`DELETE FROM orders WHERE id = $1`,
		fixture.OrderID,
	)

	_, _ = pool.Exec(
		ctx,
		`DELETE FROM wallets WHERE id = $1`,
		fixture.WalletID,
	)

	_, _ = pool.Exec(
		ctx,
		`DELETE FROM providers WHERE id = $1`,
		fixture.ProviderID,
	)

	_, _ = pool.Exec(
		ctx,
		`DELETE FROM users WHERE id = $1`,
		fixture.UserID,
	)
}
