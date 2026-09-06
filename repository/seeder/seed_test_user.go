package seeder

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/logger"
	"time"

	"go.uber.org/zap"
)

const testTelegramID = 2053084840

func (s *Seeder) SeedTestUser(ctx context.Context) error {
	const op = "seeder.SeedTestUser"

	testUsername := "test_user_phase2"

	existingUser, err := s.userRepo.FindUserByTelegramID(ctx, testTelegramID)
	if err == nil && existingUser != nil {
		logger.Logger.Info("test user already exists",
			zap.String("op", op),
			zap.Int64("telegram_id", testTelegramID),
			zap.Uint64("user_id", existingUser.ID),
		)
		return nil
	}

	user := &entity.User{
		TelegramID: testTelegramID,
		Username:   testUsername,
		Role:       entity.UserRole,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		logger.Logger.Error("failed to create test user",
			zap.String("op", op),
			zap.Error(err),
		)
		return err
	}

	logger.Logger.Info("test user created successfully",
		zap.String("op", op),
		zap.Int64("telegram_id", testTelegramID),
		zap.Uint64("user_id", user.ID),
	)

	return nil
}
