package useractivitysyncjob

type Job struct {
	userService UserService
}

func New(userService UserService) *Job {
	return &Job{userService: userService}
}

func (j *Job) Name() string {
	return "user-activity-sync"
}
