package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	PaymentsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot", Name: "payments_total",
		Help: "Total number of payments processed by method and status",
	}, []string{"method", "status"})

	WalletTransactions = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot", Name: "wallet_transactions_total",
		Help: "Total number of wallet transactions by type",
	}, []string{"type"})

	CheckoutLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "smm_bot", Name: "checkout_duration_seconds",
		Help:    "Checkout flow latency in seconds",
		Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
	}, []string{"flow_type"})
)
