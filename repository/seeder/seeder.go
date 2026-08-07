package seeder

import "context"

type Seeder interface {
	Run(ctx context.Context) error
}
