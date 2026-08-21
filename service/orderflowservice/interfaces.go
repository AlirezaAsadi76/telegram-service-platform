package orderflowservice

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/params/orderparams"
)

type Repository interface {
	Save(ctx context.Context, req orderparams.SaveOrderFlowRequest) error
	Get(ctx context.Context, req orderparams.GetOrderFlowRequest) (*orderentity.OrderFlowState, error)
	Delete(ctx context.Context, req orderparams.DeleteOrderFlowRequest) error
}
