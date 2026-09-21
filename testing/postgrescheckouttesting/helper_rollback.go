package postgrescheckouttesting

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func installFailOrderStatusTransition(t *testing.T, pool *pgxpool.Pool, orderID uint64) {
	t.Helper()

	suffix := time.Now().UnixNano()

	functionName := fmt.Sprintf(
		"test_fail_order_status_%d",
		suffix,
	)

	triggerName := fmt.Sprintf(
		"test_fail_order_trigger_%d",
		suffix,
	)

	ctx := context.Background()

	createFunctionQuery := fmt.Sprintf(`
		CREATE FUNCTION %s()
		RETURNS trigger
		LANGUAGE plpgsql
		AS $$
		BEGIN
			IF NEW.id = %d
			   AND OLD.status = 'PROCESSING'
			   AND NEW.status = 'FAILED'
			THEN
				RAISE EXCEPTION 'test induced failure during PROCESSING to FAILED';
			END IF;

			RETURN NEW;
		END;
		$$;
	`, functionName, orderID)

	if _, err := pool.Exec(
		ctx,
		createFunctionQuery,
	); err != nil {
		t.Fatalf(
			"create rollback test function: %v",
			err,
		)
	}

	createTriggerQuery := fmt.Sprintf(`
		CREATE TRIGGER %s
		BEFORE UPDATE OF status ON orders
		FOR EACH ROW
		EXECUTE FUNCTION %s();
	`, triggerName, functionName)

	if _, err := pool.Exec(
		ctx,
		createTriggerQuery,
	); err != nil {
		_, _ = pool.Exec(
			ctx,
			fmt.Sprintf(
				`DROP FUNCTION IF EXISTS %s()`,
				functionName,
			),
		)

		t.Fatalf(
			"create rollback test trigger: %v",
			err,
		)
	}

	t.Cleanup(func() {
		_, _ = pool.Exec(
			ctx,
			fmt.Sprintf(
				`DROP TRIGGER IF EXISTS %s ON orders`,
				triggerName,
			),
		)

		_, _ = pool.Exec(
			ctx,
			fmt.Sprintf(
				`DROP FUNCTION IF EXISTS %s()`,
				functionName,
			),
		)
	})
}
