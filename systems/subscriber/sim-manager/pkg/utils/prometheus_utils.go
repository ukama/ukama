package utils

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/sirupsen/logrus"
)

func PushMetrics(metricJob, metricName string, metricHelp string, metricType prometheus.ValueType, metricLabels map[string]string, metricValue float64) {
	labelDimensions := make([]string, 0, len(metricLabels))

	for key := range metricLabels {
		labelDimensions = append(labelDimensions, key)
	}

	var metric prometheus.Collector

	switch metricType {
	case prometheus.CounterValue:
		metric = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: metricName,
			Help: metricHelp,
		}, labelDimensions)
	case prometheus.GaugeValue:
		metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricName,
			Help: metricHelp,
		}, labelDimensions)
	default:
		logrus.Errorf("Unsupported metric type: %v", metricType)
		return
	}

	switch metricType {
	case prometheus.CounterValue:
		metric.(*prometheus.CounterVec).With(prometheus.Labels(metricLabels)).Add(metricValue)
	case prometheus.GaugeValue:
		metric.(*prometheus.GaugeVec).With(prometheus.Labels(metricLabels)).Set(metricValue)
	}

	if err := push.New("http://localhost:9091/", metricJob).
		Collector(metric).
		Push(); err != nil {
		logrus.Errorf("Could not push metric to Pushgateway: %s", err.Error())
	}
}
