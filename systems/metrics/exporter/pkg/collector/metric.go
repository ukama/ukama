package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/metrics/exporter/pkg"
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

func NewMetrics(name string, mtype string) *Metrics {
	m := new(Metrics)
	m.Name = name
	m.Type = MetricTypeFromString(mtype)
	return m
}

func (m *Metrics) InitializeMetric(ms pkg.MetricSchema) error {
	switch MetricTypeFromString(ms.Type) {
	case MetricGuage:
		m.collector = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:        ms.Name,
				Help:        ms.Details,
				ConstLabels: m.Labels,
			},
			ms.DynamicLabels,
		)

	case MetricCounter:
		m.collector = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        ms.Name,
				Help:        ms.Details,
				ConstLabels: m.Labels,
			},
			ms.DynamicLabels,
		)

	case MetricSummary:
		m.collector = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:        ms.Name,
				Help:        ms.Details,
				ConstLabels: m.Labels,
			},
			ms.DynamicLabels,
		)

	case MetricHistogram:
		m.collector = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        ms.Name,
				Help:        ms.Details,
				ConstLabels: m.Labels,
				Buckets:     ms.Buckets,
			},
			ms.DynamicLabels,
		)
	default:
		log.Errorf("Metric %s type %s not supported", ms.Name, ms.Type)
		return fmt.Errorf("metric %s type %s not supported", ms.Name, ms.Type)
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

func (m *Metrics) RegisterMetric(registry *prometheus.Registry) error {
	err := registry.Register(m.collector)
	if err != nil {
		log.Errorf("Failed to register metric %s. Err: %s", m.Name, err.Error())
		return err
	}
	return nil
}

func (m *Metrics) SetLabels(clabels map[string]string) {
	m.Labels = make(prometheus.Labels)

	for name, value := range clabels {
		m.Labels[name] = value
	}
}

func SetUpMetric(mc *MetricsCollector, ms pkg.MetricSchema) (*Metrics, error) {
	/* Initialize metric first */
	nm := NewMetrics(ms.Name, ms.Type)

	nm.SetLabels(ms.Labels)

	err := nm.InitializeMetric(ms)
	if err != nil {
		return nil, err
	}

	/* Add a metric */
	err = mc.AddMetrics(ms.Name, *nm)
	if err != nil {
		return nil, err
	}
	log.Infof("New metric %s added", ms.Name)

	return nm, nil

}
