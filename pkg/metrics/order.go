package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OrdersTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot", Name: "orders_total",
		Help: "Total number of orders by flow type and status",
	}, []string{"flow_type", "status"})

	ActiveOrders = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "smm_bot", Name: "active_orders",
		Help: "Current number of active orders by status",
	}, []string{"status"})

	// --- Order Flow Observability ---
	OrderFlowStateSaved = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot", Name: "order_flow_state_saved_total",
		Help: "Total number of order flow states saved",
	}, []string{"stage", "status"})

	OrderFlowStateRetrieved = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot", Name: "order_flow_state_retrieved_total",
		Help: "Total number of order flow states retrieved",
	}, []string{"found"})

	OrderFlowStateDeleted = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot", Name: "order_flow_state_deleted_total",
		Help: "Total number of order flow states deleted",
	}, []string{"reason"})

	OrderFlowDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "smm_bot", Name: "order_flow_duration_seconds",
		Help:    "Duration of order flow from start to completion",
		Buckets: []float64{10, 30, 60, 120, 300, 600},
	}, []string{"status"})
)
