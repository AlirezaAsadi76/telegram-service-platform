package postgrespaymenttesting

import (
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/paymententity"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func createTestCase(t *testing.T, pool *pgxpool.Pool, paymentStatus paymententity.PaymentStatus, orderStatus orderentity.OrderStatus) (testPayment, testOrder) {
	t.Helper()

	order := createTestOrder(t, pool, orderStatus)
	payment := createTestPayment(t, pool, order.ID, paymentStatus)

	t.Cleanup(func() {
		cleanupPaymentTestData(t, pool, payment.ID, order.ID)
	})

	return payment, order
}
