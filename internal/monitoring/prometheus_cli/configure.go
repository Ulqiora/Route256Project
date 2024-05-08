package prometheus_cli

import "github.com/prometheus/client_golang/prometheus"

func New(counters ...prometheus.Collector) *prometheus.Registry {
	reg := prometheus.NewRegistry()
	for _, counter := range counters {
		reg.MustRegister(counter)
	}
	return reg
}
