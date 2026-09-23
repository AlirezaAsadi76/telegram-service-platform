package smmcatalogsyncjob

import "context"

type ProductService interface {
	SyncSMMServices(ctx context.Context) error
}
