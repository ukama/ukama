package utils

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
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


func PushMetrics(metricJob string, metrics []struct {
    Name   string
    Type   string
    Labels map[string]string
    Value  float64
}) {
    // Extract the label keys from the first metric in the slice
    labelDimensions := make([]string, 0, len(metrics[0].Labels))
    for key := range metrics[0].Labels {
        labelDimensions = append(labelDimensions, key)
    }

    // Initialize a slice of metrics collectors for this job
    var collectors []*Metrics
    var ok bool
    if collectors, ok = metricCollectors[metricJob]; !ok {
        collectors = make([]*Metrics, 0, len(metrics))
    }

    // Iterate over each metric in the slice and add it to the collectors
    for _, metric := range metrics {
        var m *Metrics
        // Check if a collector with the same name already exists
        for _, collector := range collectors {
            if collector.Name == metric.Name {
                m = collector
                break
            }
        }
        if m == nil {
            // Create a new collector if one doesn't exist
            m = NewMetrics(metric.Name, metric.Type)
            err := m.InitializeMetric(metric.Name, MetricConfig{
                Name:   metric.Name,
                Type:   metric.Type,
                Labels: metric.Labels,
            }, labelDimensions)
            if err != nil {
                log.Errorf("Could not initialize metric collector: %s", err.Error())
                continue
            }
            collectors = append(collectors, m)
        }
        // Set the value for the current metric
        err := m.SetMetric(metric.Value, prometheus.Labels(metric.Labels))
        if err != nil {
            log.Errorf("Could not set metric value: %s", err.Error())
            continue
        }
    }

    // Save the collectors for this job
    metricCollectors[metricJob] = collectors
    // Push each collector to the Pushgateway
    for _, collector := range collectors {
        if err := push.New("http://localhost:9091/", metricJob).
            Collector(collector.collector).
            Push(); 

            err != nil {
            log.Errorf("Could not push metric to Pushgateway: %s", err.Error())
        }
    }
}
