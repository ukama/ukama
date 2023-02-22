package utils

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/sirupsen/logrus"
)

func PushMetrics(metricJob, metricName string, metricHelp string, metricLabels map[string]string, metricValue float64) {
	labelDimensions := make([]string, 0, len(metricLabels))

	for key := range metricLabels {
		labelDimensions = append(labelDimensions, key)
	}

	metric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricName,
		Help: metricHelp,
	}, labelDimensions)

	metric.With(prometheus.Labels(metricLabels)).Set(metricValue)

	if err := push.New("http://localhost:9091/", metricJob).
		Collector(metric).
		Push(); err != nil {
		logrus.Errorf("Could not push metric to Pushgateway: %s", err.Error())
	}
}
