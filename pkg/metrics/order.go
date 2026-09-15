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

	OrderFulfillmentEnqueueTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "smm_bot",
			Name:      "order_fulfillment_enqueue_total",
			Help:      "Total number of order fulfillment enqueue operations.",
		},
		[]string{"result"},
	)

	OrderFulfillmentEnqueueDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "smm_bot",
			Name:      "order_fulfillment_enqueue_duration_seconds",
			Help:      "Time spent enqueueing order fulfillment work.",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		},
	)

	OrderFlowStateSaved = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot",
		Name:      "order_flow_state_saved_total",
		Help:      "Total number of order flow states saved",
	}, []string{"stage", "status"})

	OrderFlowStateRetrieved = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot",
		Name:      "order_flow_state_retrieved_total",
		Help:      "Total number of order flow states retrieved",
	}, []string{"found"})

	OrderFlowStateDeleted = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "smm_bot",
		Name:      "order_flow_state_deleted_total",
		Help:      "Total number of order flow states deleted",
	}, []string{"reason", "status"})

	OrderFlowDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "smm_bot",
		Name:      "order_flow_duration_seconds",
		Help:      "Duration of order flow from start to completion",
		Buckets:   []float64{10, 30, 60, 120, 300, 600},
	}, []string{"status"})
)
