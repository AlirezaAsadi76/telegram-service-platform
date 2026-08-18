package productservice

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/smmentity"
)

type Repository interface {
	GetStarPlans(ctx context.Context) ([]entity.StarPackage, error)
	GetPremiumPlans(ctx context.Context) ([]entity.PremiumPlan, error)
	GetAdsPlans(ctx context.Context) ([]entity.AdsPlan, error)

	SMMMappingCreate(ctx context.Context, m *smmentity.SmmMapping) error
	SMMMappingGetAllActive(ctx context.Context) ([]smmentity.SmmMapping, error)
	SMMMappingGetByID(ctx context.Context, id int64) (*smmentity.SmmMapping, error)
	SMMMappingGetByPlatformCategory(ctx context.Context, platform smmentity.PlatformType, category smmentity.Category) ([]smmentity.SmmMapping, error)
	SMMMappingUpdate(ctx context.Context, m *smmentity.SmmMapping) error

	SMMServiceCreateOrUpdate(ctx context.Context, s smmentity.SMM) error
	SMMServiceGetAll(ctx context.Context) ([]smmentity.SMM, error)
	SMMServiceGetByD(ctx context.Context, Id int64) (*smmentity.SMM, error)
}
