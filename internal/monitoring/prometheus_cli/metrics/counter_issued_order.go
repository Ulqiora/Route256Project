package metrics

import "github.com/prometheus/client_golang/prometheus"

func NewCounterIssuedOrders() prometheus.Counter {
	return prometheus.NewCounter(prometheus.CounterOpts{
		Name: "counter_issued_orders",
		Help: "Total issued orders in service",
	})
}
