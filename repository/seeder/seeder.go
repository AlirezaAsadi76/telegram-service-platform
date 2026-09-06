package seeder

import (
	"context"
	"telegram-service-platform/repository/postgresorder"
	"telegram-service-platform/repository/postgresuser"
	"telegram-service-platform/repository/postgreswallet"
)

type Seeder struct {
	userRepo   *postgresuser.DB
	walletRepo *postgreswallet.DB
	orderRepo  *postgresorder.DB
}

func New(
	userRepo *postgresuser.DB,
	walletRepo *postgreswallet.DB,
	orderRepo *postgresorder.DB,
) *Seeder {
	return &Seeder{
		userRepo:   userRepo,
		walletRepo: walletRepo,
		orderRepo:  orderRepo,
	}
}

func (s *Seeder) SeedAll(ctx context.Context) error {

	if err := s.SeedTestUser(ctx); err != nil {
		return err
	}

	// ۲. ساخت کیف پول و تراکنش‌ها
	if err := s.SeedTestWallet(ctx); err != nil {
		return err
	}

	// ۳. ساخت سفارشات
	if err := s.SeedTestOrders(ctx); err != nil {
		return err
	}

	return nil
}
