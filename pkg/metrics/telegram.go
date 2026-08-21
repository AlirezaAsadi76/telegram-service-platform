package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TelegramUpdates = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot", Name: "telegram_updates_total",
		Help: "Total number of Telegram updates received by type",
	}, []string{"type", "status"})

	TelegramHandlerDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "smm_bot", Name: "telegram_handler_duration_seconds",
		Help:    "Duration of Telegram handler execution",
		Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	}, []string{"update_type", "handler_hint"})

	TelegramPanics = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "smm_bot", Name: "telegram_panics_total",
		Help: "Total number of panics recovered in Telegram handlers",
	})
)
