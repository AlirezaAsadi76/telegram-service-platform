package orderflowtesting

import (
	"context"

	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/params/orderparams"
)

type fakeRepository struct {
	saveReq *orderparams.SaveOrderFlowRequest

	getState *orderentity.OrderFlowState
	getErr   error

	deleteCalls int
	deleteErr   error
}

func (f *fakeRepository) Save(
	_ context.Context,
	req orderparams.SaveOrderFlowRequest,
) error {
	copyReq := req
	copyReq.State = req.State

	f.saveReq = &copyReq

	return nil
}

func (f *fakeRepository) Get(
	_ context.Context,
	_ orderparams.GetOrderFlowRequest,
) (*orderentity.OrderFlowState, error) {
	return f.getState, f.getErr
}

func (f *fakeRepository) Delete(
	_ context.Context,
	_ orderparams.DeleteOrderFlowRequest,
) error {
	f.deleteCalls++

	return f.deleteErr
}
