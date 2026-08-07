package scheduler

import (
	"context"
)

func (s *Scheduler) Shutdown(ctx context.Context) error {

	return s.engine.Shutdown()

}
