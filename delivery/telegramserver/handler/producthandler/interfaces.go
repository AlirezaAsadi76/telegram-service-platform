package producthandler

import (
	"context"
	"telegram-service-platform/params/productparams"
)

type ProductService interface {
	GetStarPlans(ctx context.Context) (productparams.GetStarPlansResponse, error)
	GetPremiumPlans(ctx context.Context) (productparams.GetPremiumPlansResponse, error)
}
