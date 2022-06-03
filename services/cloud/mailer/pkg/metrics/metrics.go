package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func EmailSentSuccessfulRequestMetric() {
	go func() {
		opsSuccessProcessed.Inc()
	}()
}

func EmailSentFailureRequestMetric() {
	go func() {
		opsFailedProcessed.Inc()
	}()
}

var (
	opsSuccessProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "mailer_email_processed",
		Help:        "The total emails processed",
		ConstLabels: map[string]string{"status": "succeeded"},
	})
)

var (
	opsFailedProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "mailer_email_processed",
		Help:        "The total emails processed",
		ConstLabels: map[string]string{"status": "failed"},
	})
)
