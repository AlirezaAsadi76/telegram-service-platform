package metrics

import "time"

func RecordWorkerRun(jobName string, fn func()) {
	start := time.Now()
	fn()
	WorkerDuration.WithLabelValues(jobName).Observe(time.Since(start).Seconds())
}
