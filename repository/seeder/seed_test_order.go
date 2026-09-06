package seeder

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/logger"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// SeedTestOrders سفارشات تستی با وضعیت‌های مختلف ایجاد می‌کند
func (s *Seeder) SeedTestOrders(ctx context.Context) error {
	const op = "seeder.SeedTestOrders"

	// دریافت کاربر تستی
	user, err := s.userRepo.FindUserByTelegramID(ctx, testTelegramID)
	if err != nil {
		logger.Logger.Error("test user not found",
			zap.String("op", op),
			zap.Error(err),
		)
		return err
	}

	orders := []orderentity.Order{
		{
			UserID:      user.ID,
			ProductType: productentity.ProductTypeSMM,
			ProductID:   1,
			Quantity:    1000,
			Amount:      entity.Amount(decimal.NewFromInt(15000)),
			Currency:    "TOMAN",
			Status:      orderentity.OrderStatusPending,
			CreatedAt:   time.Now().Add(-2 * time.Hour),
			UpdatedAt:   time.Now().Add(-2 * time.Hour),
		},
		{
			UserID:      user.ID,
			ProductType: productentity.ProductTypeSMM,
			ProductID:   2,
			Quantity:    500,
			Amount:      entity.Amount(decimal.NewFromInt(7500)),
			Currency:    "TOMAN",
			Status:      orderentity.OrderStatusPaid,
			CreatedAt:   time.Now().Add(-1 * time.Hour),
			UpdatedAt:   time.Now().Add(-1 * time.Hour),
		},
		{
			UserID:      user.ID,
			ProductType: productentity.ProductTypeSMM,
			ProductID:   3,
			Quantity:    2000,
			Amount:      entity.Amount(decimal.NewFromInt(30000)),
			Currency:    "TOMAN",
			Status:      orderentity.OrderStatusProcessing,
			CreatedAt:   time.Now().Add(-30 * time.Minute),
			UpdatedAt:   time.Now().Add(-30 * time.Minute),
		},
		{
			UserID:      user.ID,
			ProductType: productentity.ProductTypeSMM,
			ProductID:   4,
			Quantity:    1500,
			Amount:      entity.Amount(decimal.NewFromInt(22500)),
			Currency:    "TOMAN",
			Status:      orderentity.OrderStatusSuccess,
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now().Add(-23 * time.Hour),
		},
		{
			UserID:      user.ID,
			ProductType: productentity.ProductTypeSMM,
			ProductID:   5,
			Quantity:    800,
			Amount:      entity.Amount(decimal.NewFromInt(12000)),
			Currency:    "TOMAN",
			Status:      orderentity.OrderStatusFailed,
			CreatedAt:   time.Now().Add(-48 * time.Hour),
			UpdatedAt:   time.Now().Add(-47 * time.Hour),
		},
	}

	for _, order := range orders {
		err := s.orderRepo.Create(ctx, &order)
		if err != nil {
			logger.Logger.Error("failed to create test order",
				zap.String("op", op),
				zap.Error(err),
			)
			return err
		}
	}

	logger.Logger.Info("test orders created successfully",
		zap.String("op", op),
		zap.Int("count", len(orders)),
	)

	return nil
}
