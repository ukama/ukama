package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func RecordIpRequestSuccessMetric() {
	go func() {
		opsSuccessProcessed.Inc()
	}()
}

func RecordIpRequestFailureMetric() {
	go func() {
		opsFailedProcessed.Inc()
	}()
}

func RecordSetIpMetric() {
	go func() {
		opsSetIPRequest.Inc()
	}()
}

const getIpRequests = "net_get_ip_requests_total"
const getIpRequestsDescr = "The total number of get IP requests"

var (
	opsSuccessProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name:        getIpRequests,
		Help:        getIpRequestsDescr,
		ConstLabels: map[string]string{"status": "succeeded"},
	})
)

var (
	opsFailedProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name:        getIpRequests,
		Help:        getIpRequestsDescr,
		ConstLabels: map[string]string{"status": "failed"},
	})
)

var (
	opsSetIPRequest = promauto.NewCounter(prometheus.CounterOpts{
		Name: "net_set_ip_requests_total",
		Help: "The total number of set ip requests",
	})
)
