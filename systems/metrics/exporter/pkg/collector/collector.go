package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/metrics/exporter/pkg"
)

type Metrics struct {
	Name      string
	gauge     *prometheus.GaugeVec
	counter   *prometheus.CounterVec
	summary   *prometheus.SummaryVec
	histogram *prometheus.HistogramVec
}

type MetricsCollector struct {
	MetricsMap map[string]Metrics
	Config     map[string]pkg.KPIConfig
}

func NewMetricsCollector(config []pkg.KPIConfig) *MetricsCollector {
	m := new(MetricsCollector)
	m.MetricsMap = make(map[string]Metrics)
	m.Config = make(map[string]pkg.KPIConfig, len(m.Config))
	for _, c := range config {
		m.Config[c.Event] = c
	}
	return m
}

func (c *MetricsCollector) GetConfigForEvent(event string) (*pkg.KPIConfig, error) {

	cfg, ok := c.Config[event]
	if !ok {
		log.Errorf("Event %s not expected by exporter service.", event)
		return nil, fmt.Errorf("event %s not supported", event)
	}

	return &cfg, nil
}

func (c *MetricsCollector) GetMetric(name string) (*Metrics, error) {

	m, ok := c.MetricsMap[name]
	if !ok {
		log.Errorf("Metric %s doesn't exist", name)
		return nil, fmt.Errorf("Metric %s doesn't exist", name)
	}

	return &m, nil
}

func (c *MetricsCollector) AddMetrics(name string, m Metrics) error {
	_, ok := c.MetricsMap[name]
	if !ok {
		c.MetricsMap[name] = m
	} else {
		log.Errorf("Metric %s already exist", name)
		return fmt.Errorf("Metric %s already exist", name)
	}
	return nil
}

func (m *Metrics) InitializeMetric(name string, config pkg.KPIConfig, customLables []string) {
	switch config.Type {
	case pkg.MetricGuage:
		m.gauge = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:        config.Name,
				Help:        config.Details,
				ConstLabels: config.Labels,
			},
			customLables,
		)
	case pkg.MetricCounter:
		m.counter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        config.Name,
				Help:        config.Details,
				ConstLabels: config.Labels,
			},
			customLables,
		)
	case pkg.MetricSummary:
		m.summary = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:        config.Name,
				Help:        config.Details,
				ConstLabels: config.Labels,
			},
			customLables,
		)
	case pkg.MetricHistogram:
		m.histogram = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        config.Name,
				Help:        config.Details,
				ConstLabels: config.Labels,
			},
			customLables,
		)
	}
}
