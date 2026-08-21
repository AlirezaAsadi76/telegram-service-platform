package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	SMMProviderRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot", Name: "smm_provider_requests_total",
		Help: "Total SMM provider API requests by provider and status",
	}, []string{"provider", "status"})

	SMMProviderLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "smm_bot", Name: "smm_provider_duration_seconds",
		Help:    "SMM provider API latency in seconds",
		Buckets: []float64{0.1, 0.25, 0.5, 1, 2, 5, 10, 30},
	}, []string{"provider"})

	SMMServiceSyncTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot", Name: "smm_service_sync_total",
		Help: "Total number of SMM services synced from provider",
	}, []string{"status"})

	SMMServiceSyncDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "smm_bot", Name: "smm_service_sync_duration_seconds",
		Help:    "Duration of SMM service sync job in seconds",
		Buckets: []float64{1, 2, 5, 10, 30, 60, 120},
	})

	SMMCatalogSize = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "smm_bot", Name: "smm_catalog_size",
		Help: "Current number of active SMM service mappings",
	})
)
