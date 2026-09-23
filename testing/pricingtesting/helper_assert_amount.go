package pricingtesting

import (
	"telegram-service-platform/entity"
	"testing"

	"github.com/shopspring/decimal"
)

func assertAmount(t *testing.T, got entity.Amount, want float64) {
	t.Helper()

	expected := entity.Amount(
		decimal.NewFromFloat(want),
	)

	if !got.Equal(expected) {
		t.Fatalf(
			"expected amount %s, got %s",
			expected.String(),
			got.String(),
		)
	}
}
