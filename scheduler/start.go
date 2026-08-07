package scheduler

func (s *Scheduler) Start() {

	s.engine.Start()

	go func() {

		<-s.ctx.Done()

		s.engine.Shutdown()

	}()

}
