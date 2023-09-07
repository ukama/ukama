package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func RecordSuccessfulRequestMetric() {
	go func() {
		opsSuccessProcessed.Inc()
	}()
}

func RecordFailedRequestMetric() {
	go func() {
		opsFailedProcessed.Inc()
	}()
}

var (
	opsSuccessProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "foo_requests_total",
		Help:        "The total number requests",
		ConstLabels: map[string]string{"status": "succeeded"},
	})
)

var (
	opsFailedProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "foo_requests_total",
		Help:        "The total number requests",
		ConstLabels: map[string]string{"status": "failed"},
	})
)
