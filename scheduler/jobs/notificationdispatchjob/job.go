package notificationdispatchjob

import (
	"sync"

	"telegram-service-platform/service/notificationservice"
)

type Job struct {
	notificationService notificationservice.Service
	redis               RedisRepository
	bot                 TelegramSender
	mutex               sync.Mutex
	config              Config
}

func New(
	notificationService notificationservice.Service,
	redis RedisRepository,
	bot TelegramSender,
	config Config,
) *Job {
	return &Job{
		notificationService: notificationService,
		redis:               redis,
		bot:                 bot,
		config:              config,
	}
}

func (j *Job) Name() string {
	return "notification-dispatch"
}
