package userservice

type Service struct {
	repository      UserRepository
	activityTracker ActivityTrackerRepository
}

func New(repository UserRepository, activityTracker ActivityTrackerRepository) *Service {
	return &Service{
		repository:      repository,
		activityTracker: activityTracker,
	}
}
