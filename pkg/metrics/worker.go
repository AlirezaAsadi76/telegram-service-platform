package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	WorkerRuns = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot", Name: "worker_runs_total",
		Help: "Total number of worker executions by job name and status",
	}, []string{"job_name", "status"})

	WorkerDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "smm_bot", Name: "worker_duration_seconds",
		Help:    "Worker execution duration in seconds",
		Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2, 5, 10, 30, 60},
	}, []string{"job_name"})

	NotificationsSent = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot", Name: "notifications_total",
		Help: "Total notifications sent by status",
	}, []string{"status"})

	NotificationQueueDepth = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "smm_bot", Name: "notification_queue_depth",
		Help: "Current depth of notification Redis queue",
	})
)
