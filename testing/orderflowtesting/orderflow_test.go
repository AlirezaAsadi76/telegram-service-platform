package orderflowtesting

import (
	"context"
	"testing"
	"time"

	"telegram-service-platform/entity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/service/orderflowservice"
)

func TestService_SaveOrderFlow_StartsTTLOnlyOnce(
	t *testing.T,
) {
	repo := &fakeRepository{}

	service := orderflowservice.New(
		repo,
		orderflowservice.Config{
			OrderTTL: 15 * time.Minute,
		},
	)

	now := time.Now()

	err := service.SaveOrderFlow(
		context.Background(),
		orderparams.SaveOrderFlowRequest{
			TelegramID: entity.TelegramId(100),
			State: orderentity.OrderFlowState{
				Stage:     orderentity.OrderFlowStageWaitingForQuantity,
				ServiceID: 55,
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.saveReq == nil {
		t.Fatal("expected repository Save to be called")
	}

	firstExpiresAt := repo.saveReq.State.ExpiresAt
	expiresAt := time.Unix(firstExpiresAt, 0)

	if expiresAt.Before(
		now.Add(14 * time.Minute),
	) {
		t.Fatalf(
			"expected expiry around 15 minutes from now, got %v",
			expiresAt,
		)
	}

	firstTTL := repo.saveReq.TTLMins

	if firstTTL <= 0 ||
		firstTTL > 15*time.Minute {
		t.Fatalf(
			"expected remaining TTL within 15 minutes, got %v",
			firstTTL,
		)
	}

	err = service.SaveOrderFlow(
		context.Background(),
		orderparams.SaveOrderFlowRequest{
			TelegramID: entity.TelegramId(100),
			State: orderentity.OrderFlowState{
				Stage:     orderentity.OrderFlowStageWaitingForLink,
				ServiceID: 55,
				ExpiresAt: firstExpiresAt,
			},
		},
	)
	if err != nil {
		t.Fatalf(
			"unexpected error on second save: %v",
			err,
		)
	}

	if repo.saveReq.State.ExpiresAt != firstExpiresAt {
		t.Fatalf(
			"expected expiry to remain %d, got %d",
			firstExpiresAt,
			repo.saveReq.State.ExpiresAt,
		)
	}

	if repo.saveReq.TTLMins >= firstTTL {
		t.Fatalf(
			"expected remaining TTL to decrease, first=%v current=%v",
			firstTTL,
			repo.saveReq.TTLMins,
		)
	}
}

func TestService_GetOrderFlow_ExpiredStateIsInvalidated(
	t *testing.T,
) {
	repo := &fakeRepository{
		getState: &orderentity.OrderFlowState{
			Stage:     orderentity.OrderFlowStageConfirming,
			ServiceID: 55,
			ExpiresAt: time.Now().Add(-1 * time.Minute).Unix(),
		},
	}

	service := orderflowservice.New(
		repo,
		orderflowservice.Config{
			OrderTTL: 15 * time.Minute,
		},
	)

	state, err := service.GetOrderFlow(
		context.Background(),
		orderparams.GetOrderFlowRequest{
			TelegramID: entity.TelegramId(100),
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if state != nil {
		t.Fatalf(
			"expected expired flow to be unavailable, got %#v",
			state,
		)
	}

	if repo.deleteCalls != 1 {
		t.Fatalf(
			"expected expired flow to be deleted once, got %d",
			repo.deleteCalls,
		)
	}
}
