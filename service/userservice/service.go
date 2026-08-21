package userservice

type Service struct {
	repository      UserRepository
	activityTracker ActivityTrackerRepository
	walletSvc       WalletService
}

func New(walletSvc WalletService, repository UserRepository, activityTracker ActivityTrackerRepository) *Service {
	return &Service{
		repository:      repository,
		activityTracker: activityTracker,
		walletSvc:       walletSvc,
	}
}
