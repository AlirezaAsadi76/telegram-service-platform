package producthandler

import (
	"context"
	"telegram-service-platform/params"
)

type ProductService interface {
	GetStarPlans(ctx context.Context) (params.GetStarPlansResponse, error)
	GetPremiumPlans(ctx context.Context) (params.GetPremiumPlansResponse, error)
}
