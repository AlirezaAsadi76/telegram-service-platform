package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	PaymentIntentResult = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "smm_bot",
			Name:      "payment_intent_total",
			Help:      "Total number of payment intent operations by method and result.",
		},
		[]string{"method", "result"},
	)

	PaymentIntentDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "smm_bot",
			Name:      "payment_intent_duration_seconds",
			Help:      "Time spent creating payment intents.",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		},
	)

	PaymentInitiationResult = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "smm_bot",
			Name:      "payment_initiation_total",
			Help:      "Total number of payment initiation operations.",
		},
		[]string{"method", "result"},
	)

	PaymentInitiationDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "smm_bot",
			Name:      "payment_initiation_duration_seconds",
			Help:      "Time spent initiating external payments.",
			Buckets:   []float64{0.05, 0.1, 0.25, 0.5, 1, 2, 5, 10},
		},
	)
)
