package scheduler

import "context"

func (s *Scheduler) Start(ctx context.Context) {
	s.ctx = ctx
	s.engine.Start()

	go func() {

		<-s.ctx.Done()

		s.engine.Shutdown()

	}()

}
