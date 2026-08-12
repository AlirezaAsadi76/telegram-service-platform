package checkoutservice

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/params/walletparam"
)

type WalletService interface {
	Credit(ctx context.Context, req walletparam.CreditRequest) (*walletparam.CreditResponse, error)
	Debit(ctx context.Context, req walletparam.DebitRequest) (*walletparam.DebitResponse, error)
	GetBalance(ctx context.Context, req walletparam.GetBalanceRequest) (*walletparam.GetBalanceResponse, error)
}

type PaymentService interface {
	Create(ctx context.Context, req paymentparams.CreateRequest) (*paymentparams.CreateResponse, error)
	Verify(ctx context.Context, req paymentparams.VerifyRequest) (*paymentparams.VerifyResponse, error)
	GetById(ctx context.Context, paymentID uint64) (*paymententity.Payment, error)
}

type OrderService interface {
	Create(ctx context.Context, req orderparams.CreateRequest) (*orderparams.CreateResponse, error)
	UpdateStatus(ctx context.Context, req orderparams.UpdateStatusRequest) error
	GetById(ctx context.Context, orderId uint64) (*orderentity.Order, error)
}

type SMMProviderService interface {
	FulfillOrder(ctx context.Context, order *orderentity.Order) error
}

type Messenger interface {
	SendToUser(ctx context.Context, userID int64, message string) error
	SendToAdminChannel(ctx context.Context, message string) error
}
