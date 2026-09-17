package justanotherpanel

import (
	"context"
	"strings"

	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (a *Adapter) GetOrderStatus(ctx context.Context, externalOrderID string) (orderentity.OrderStatus, error) {
	const Op = "justanotherpanel.GetOrderStatus"

	response, err := a.Status(ctx, externalOrderID)
	if err != nil {
		return "", richerror.New(Op, err).
			WithKind(richerror.KindExternalAPI).
			WithCode(richerror.CodeSMMProviderRequestFailed).
			WithMessage(msgerror.SMMProviderRequestFailed)
	}

	status := strings.ToLower(strings.TrimSpace(response.Status.Status))

	switch StatusType(status) {
	case statusCompleted:
		return orderentity.OrderStatusCompleted, nil

	case statusInProcessing,
		statusPending,
		statusProcessing:
		return orderentity.OrderStatusProcessing, nil

	case statusCancelled,
		statusFailed:
		return orderentity.OrderStatusFailed, nil

	case "partial":
		return orderentity.OrderStatusProcessing, nil

	default:
		return "", richerror.New(Op, nil).
			WithKind(richerror.KindExternalAPI).
			WithCode(richerror.CodeSMMProviderInvalidResponse).
			WithMessage(msgerror.SMMProviderInvalidResponse)
	}
}
