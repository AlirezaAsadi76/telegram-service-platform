package producttesting

import (
	"context"
	"telegram-service-platform/params/smmparams"
	"telegram-service-platform/service/productservice"
)

var _ productservice.SMMAdapterInterface = (*fakeSMMAdapter)(nil)

type fakeSMMAdapter struct {
	response smmparams.GetAllServicesResponse
	err      error

	calls int
}

func (f *fakeSMMAdapter) AllServices(
	_ context.Context,
) (smmparams.GetAllServicesResponse, error) {
	f.calls++

	return f.response, f.err
}
