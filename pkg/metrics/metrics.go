package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
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

	// ✅ تعریف واحد و نهایی برای OrdersTotal (حذف OrdersCreated تکراری)
	OrdersTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot",
		Name:      "orders_total",
		Help:      "Total number of orders by flow type and status",
	}, []string{"flow_type", "status"})

	// SMMServiceSyncTotal counts how many services were synced from JAP.
	SMMServiceSyncTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot",
		Name:      "smm_service_sync_total",
		Help:      "Total number of SMM services synced from provider",
	}, []string{"status"})

	// SMMServiceSyncDuration tracks the time taken for a full sync.
	SMMServiceSyncDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "smm_bot",
		Name:      "smm_service_sync_duration_seconds",
		Help:      "Duration of SMM service sync job in seconds",
		Buckets:   []float64{1, 2, 5, 10, 30, 60, 120},
	})

	// SMMCatalogSize is a gauge showing how many active mappings exist.
	SMMCatalogSize = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "smm_bot",
		Name:      "smm_catalog_size",
		Help:      "Current number of active SMM service mappings",
	})

	// Telegram Bot Observability
	TelegramUpdates = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "smm_bot",
			Name:      "telegram_updates_total",
			Help:      "Total number of Telegram updates received by type",
		},
		[]string{"type", "status"}, // type: message, callback, command | status: success, panic, error
	)

	TelegramHandlerDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "smm_bot",
			Name:      "telegram_handler_duration_seconds",
			Help:      "Duration of Telegram handler execution",
			Buckets:   []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"update_type", "handler_hint"},
	)

	TelegramPanics = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "smm_bot",
			Name:      "telegram_panics_total",
			Help:      "Total number of panics recovered in Telegram handlers",
		},
	)

	// Order Flow Observability
	OrderFlowStateSaved = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "smm_bot",
			Name:      "order_flow_state_saved_total",
			Help:      "Total number of order flow states saved",
		},
		[]string{"stage"},
	)

	OrderFlowStateRetrieved = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "smm_bot",
			Name:      "order_flow_state_retrieved_total",
			Help:      "Total number of order flow states retrieved",
		},
		[]string{"found"},
	)

	OrderFlowStateDeleted = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "smm_bot",
			Name:      "order_flow_state_deleted_total",
			Help:      "Total number of order flow states deleted",
		},
		[]string{"reason"}, // "completed", "abandoned", "expired"
	)

	OrderFlowDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "smm_bot",
			Name:      "order_flow_duration_seconds",
			Help:      "Duration of order flow from start to completion",
			Buckets:   []float64{10, 30, 60, 120, 300, 600}, // 10s, 30s, 1m, 2m, 5m, 10m
		},
		[]string{"status"}, // "completed", "abandoned"
	)
)
