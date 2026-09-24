package orderfulfillertesting

import (
	"context"

	"telegram-service-platform/entity/smmentity"
)

type fakeProductService struct {
	service *smmentity.SMM
	err     error
	calls   int

	lastMappingID int64
}

func (f *fakeProductService) GetSMMServiceByMappingID(
	_ context.Context,
	mappingID int64,
) (*smmentity.SMM, error) {
	f.calls++
	f.lastMappingID = mappingID

	return f.service, f.err
}
