package orderfulfillmentservice

type Service struct {
	queue  QueueRepository
	config Config
}

func New(queue QueueRepository, config Config) *Service {
	return &Service{
		queue:  queue,
		config: config,
	}
}
