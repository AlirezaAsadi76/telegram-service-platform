package scheduler

import (
	"context"

	"github.com/go-co-op/gocron/v2"

	"telegram-service-platform/scheduler/jobs"
)

type Scheduler struct {
	engine gocron.Scheduler
	jobs   []jobs.Job
	config Config
	ctx    context.Context
}

func New(config Config, jobList ...jobs.Job) (*Scheduler, error) {

	engine, err := gocron.NewScheduler()

	if err != nil {
		return nil, err
	}

	return &Scheduler{
		engine: engine,
		jobs:   jobList,
		config: config,
	}, nil
}
