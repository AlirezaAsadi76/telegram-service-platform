package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// OrdersCreated — تعداد سفارش‌های ایجاد شده (label: flow_type, status)
	OrdersCreated = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot",
		Name:      "orders_total",
		Help:      "Total number of orders created by flow type and status",
	}, []string{"flow_type", "status"})

	// PaymentsProcessed — تعداد پرداخت‌ها (label: method, status)
	PaymentsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot",
		Name:      "payments_total",
		Help:      "Total number of payments processed by method and status",
	}, []string{"method", "status"})

	// WalletTransactions — تراکنش‌های کیف پول (label: type)
	WalletTransactions = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot",
		Name:      "wallet_transactions_total",
		Help:      "Total number of wallet transactions by type",
	}, []string{"type"})

	// ActiveOrders — تعداد سفارش‌های فعال (gauge)
	ActiveOrders = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "smm_bot",
		Name:      "active_orders",
		Help:      "Current number of active orders by status",
	}, []string{"status"})

	// CheckoutLatency — زمان پاسخ Checkout (histogram)
	CheckoutLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "smm_bot",
		Name:      "checkout_duration_seconds",
		Help:      "Checkout flow latency in seconds",
		Buckets:   []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
	}, []string{"flow_type"})

	// WorkerRuns — تعداد اجرای Worker (label: job_name, status)
	WorkerRuns = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot",
		Name:      "worker_runs_total",
		Help:      "Total number of worker executions by job name and status",
	}, []string{"job_name", "status"})

	// WorkerDuration — زمان اجرای Worker (histogram)
	WorkerDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "smm_bot",
		Name:      "worker_duration_seconds",
		Help:      "Worker execution duration in seconds",
		Buckets:   []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2, 5, 10, 30, 60},
	}, []string{"job_name"})

	// SMMProviderRequests — تعداد درخواست به SMM Provider (label: provider, status)
	SMMProviderRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot",
		Name:      "smm_provider_requests_total",
		Help:      "Total SMM provider API requests by provider and status",
	}, []string{"provider", "status"})

	// SMMProviderLatency — زمان پاسخ SMM Provider (histogram)
	SMMProviderLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "smm_bot",
		Name:      "smm_provider_duration_seconds",
		Help:      "SMM provider API latency in seconds",
		Buckets:   []float64{0.1, 0.25, 0.5, 1, 2, 5, 10, 30},
	}, []string{"provider"})

	// NotificationsSent — تعداد نوتیفیکیشن (label: status)
	NotificationsSent = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot",
		Name:      "notifications_total",
		Help:      "Total notifications sent by status",
	}, []string{"status"})

	// NotificationQueueDepth — عمق صف Redis (gauge)
	NotificationQueueDepth = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "smm_bot",
		Name:      "notification_queue_depth",
		Help:      "Current depth of notification Redis queue",
	})

	// SMMServiceSyncTotal counts how many services were synced from JAP.
	SMMServiceSyncTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot",
		Name:      "smm_service_sync_total",
		Help:      "Total number of SMM services synced from provider",
	}, []string{"status"})

	// SMMServiceSyncDuration tracks the time taken for a full sync.
	SMMServiceSyncDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "smm_bot",
		Name:      "smm_service_sync_duration_seconds",
		Help:      "Duration of SMM service sync job in seconds",
		Buckets:   []float64{1, 2, 5, 10, 30, 60, 120},
	})

	// SMMCatalogSize is a gauge showing how many active mappings exist.
	SMMCatalogSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "smm_bot",
		Name:      "smm_catalog_size",
		Help:      "Current number of active SMM service mappings",
	})
)
