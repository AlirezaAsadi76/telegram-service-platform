package scheduler

import "context"

func (s *Scheduler) Start(ctx context.Context) {

	s.engine.Start()

	go func() {

		<-ctx.Done()

		s.engine.Shutdown()

	}()

}
