package utils

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type MetricType string

const (
	MetricUnkown    MetricType = "unavailable"
	MetricGuage     MetricType = "guage"
	MetricCounter   MetricType = "counter"
	MetricHistogram MetricType = "histogram"
	MetricSummary   MetricType = "summary"
)

type Metrics struct {
	Name      string
	Type      MetricType
	collector prometheus.Collector
	Labels    prometheus.Labels
}
type MetricConfig struct {
	Name    string
	Event   string
	Type    string
	Units   string
	Labels  map[string]string
	Details string
	Buckets []float64
}

func MetricTypeFromString(s string) MetricType {
	switch s {
	case "gauge":
		return MetricGuage
	case "counter":
		return MetricCounter
	case "histogram":
		return MetricHistogram
	case "summary":
		return MetricSummary
	default:
		return MetricUnkown
	}
}

var metricCollectors = make(map[string][]*Metrics)

func NewMetrics(name string, mtype string) *Metrics {
	m := new(Metrics)
	m.Name = name
	m.Type = MetricTypeFromString(mtype)
	return m
}

func (m *Metrics) InitializeMetric(name string, config MetricConfig, customLables []string) error {
	switch MetricTypeFromString(config.Type) {
	case MetricGuage:
		m.collector = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:        config.Name,
				Help:        config.Details,
				ConstLabels: m.Labels,
			},
			customLables,
		)

	case MetricCounter:
		m.collector = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        config.Name,
				Help:        config.Details,
				ConstLabels: m.Labels,
			},
			customLables,
		)

	case MetricSummary:
		m.collector = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:        config.Name,
				Help:        config.Details,
				ConstLabels: m.Labels,
			},
			customLables,
		)

	case MetricHistogram:
		m.collector = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        config.Name,
				Help:        config.Details,
				ConstLabels: m.Labels,
				Buckets:     config.Buckets,
			},
			customLables,
		)
	default:
		log.Errorf("Metric %s type %s not supported", config.Name, config.Type)
		return fmt.Errorf("metric %s type %s not supported", config.Name, config.Type)
	}

	return nil
}

func (m *Metrics) SetMetric(value float64, labels prometheus.Labels) error {
	switch met := m.collector.(type) {
	case *prometheus.GaugeVec:
		met.With(labels).Set(value)
	case *prometheus.CounterVec:
		met.With(labels).Inc()
	case *prometheus.SummaryVec:
		met.With(labels).Observe(value)
	case *prometheus.HistogramVec:
		met.With(labels).Observe(value)
	default:
		return fmt.Errorf("metric type %s not supported", m.Type)
	}
	return nil
}
func PushMetrics(metricJob, metricName string, metricType string, metricLabels map[string]string, metricValue float64) {
	labelDimensions := make([]string, 0, len(metricLabels))

	for key := range metricLabels {
		labelDimensions = append(labelDimensions, key)
	}

	if collectors, ok := metricCollectors[metricJob]; ok {
		for _, collector := range collectors {
			if collector.Name == metricName {
				// Metric collector exists, set metric value
				collector.SetMetric(metricValue, prometheus.Labels(metricLabels))
				return
			}
		}
	}

	metric := NewMetrics(metricName, metricType)
	err := metric.InitializeMetric(metricName, MetricConfig{
		Name:   metricName,
		Type:   string(metricType),
		Labels: metricLabels,
	}, labelDimensions)
	if err != nil {
		logrus.Errorf("Could not initialize metric collector: %s", err.Error())
		return
	}

	metricCollectors[metricJob] = append(metricCollectors[metricJob], metric)

	err = metric.SetMetric(metricValue, prometheus.Labels(metricLabels))
	if err != nil {
		logrus.Errorf("Could not set metric value: %s", err.Error())
		return
	}

	if err := push.New("http://localhost:9091/", metricJob).
		Collector(metric.collector).
		Push(); err != nil {
		logrus.Errorf("Could not push metric to Pushgateway: %s", err.Error())
	}
}
