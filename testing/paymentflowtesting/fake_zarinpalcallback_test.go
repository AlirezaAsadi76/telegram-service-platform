package paymentflowtesting

import (
	"telegram-service-platform/delivery/httpserver/paymenthandler"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/repository/postgrespayment"
	"telegram-service-platform/service/paymentservice"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

func newZarinpalCallbackFlow(t *testing.T, provider *flowTestProvider, validator *flowTestValidator) (
	*echo.Echo, *pgxpool.Pool) {
	t.Helper()

	pool := newTestPool(t)

	transactionProvider := postgres.NewTransactionProvider(pool)

	paymentRepo := postgrespayment.NewWithExecutor(
		pool,
		transactionProvider,
	)

	paymentService := paymentservice.New(
		paymentRepo,
		paymentRepo,
		provider,
		nil,
	)

	handler := paymenthandler.New(
		paymentService,
		validator,
	)

	e := echo.New()

	handler.SetRoutes(e)

	return e, pool
}
