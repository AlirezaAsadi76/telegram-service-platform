package smmcatalogsyncjobtesting

import (
	"context"
	"telegram-service-platform/scheduler/jobs/smmcatalogsyncjob"
)

var _ smmcatalogsyncjob.ProductService = (*fakeProductService)(nil)

type fakeProductService struct {
	calls int
	err   error
}

func (f *fakeProductService) SyncSMMServices(
	_ context.Context,
) error {
	f.calls++

	return f.err
}
