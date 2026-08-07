package scheduler

func (s *Scheduler) Shutdown() error {

	return s.engine.Shutdown()

}
