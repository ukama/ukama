package queue

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func MessageProcessedMetric() {
	go func() {
		opsSuccessProcessedSuccess.Inc()
	}()
}

func MessageProcessFailedMetric() {
	go func() {
		opsSuccessProcessedFailure.Inc()
	}()
}

const messageProcessedSuccess = "queue_message_processed"
const messageProcessedDescr = "The total number of processed messages"

var (
	opsSuccessProcessedSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name:        messageProcessedSuccess,
		Help:        messageProcessedDescr,
		ConstLabels: map[string]string{"status": "succeeded"},
	})
)

var (
	opsSuccessProcessedFailure = promauto.NewCounter(prometheus.CounterOpts{
		Name:        messageProcessedSuccess,
		Help:        messageProcessedDescr,
		ConstLabels: map[string]string{"status": "failed"},
	})
)
