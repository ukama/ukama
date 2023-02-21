package collector

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	pc "github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/metrics/exporter/pkg"
)

type Metrics struct {
	Name      string
	Type      pkg.MetricType
	collector prometheus.Collector
	Labels    prometheus.Labels
}

type MetricsCollector struct {
	MetricsMap map[string]Metrics
	Config     map[string]pkg.KPIConfig
	registry   *prometheus.Registry
}

func NewMetricsCollector(config []pkg.KPIConfig, metrics *config.Metrics) *MetricsCollector {
	m := new(MetricsCollector)
	m.MetricsMap = make(map[string]Metrics)
	m.Config = make(map[string]pkg.KPIConfig, len(m.Config))
	m.registry = prometheus.NewRegistry()
	m.registry.MustRegister(pc.NewGoCollector(), pc.NewProcessCollector(pc.ProcessCollectorOpts{}))

	for _, c := range config {
		m.Config[c.Event] = c
	}

	m.StartMetricServer(metrics)

	return m
}

func (c *MetricsCollector) StartMetricServer(metrics *config.Metrics) {
	go func() {
		handler := promhttp.HandlerFor(c.registry, promhttp.HandlerOpts{})
		http.Handle("/metrics", handler)
		log.Infof("Starting metrics server on port %d", metrics.Port)
		err := http.ListenAndServe(fmt.Sprintf(":%d", metrics.Port), nil)
		if err != nil {
			log.Fatalf("Error starting metrics server: %s", err.Error())
		}
	}()
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
		return nil, fmt.Errorf("metric %s doesn't exist", name)
	}

	return &m, nil
}

func (c *MetricsCollector) AddMetrics(name string, m Metrics) error {
	_, ok := c.MetricsMap[name]
	if !ok {
		c.MetricsMap[name] = m
		err := m.RegisterMetric(c.registry)
		if err != nil {
			log.Errorf("Metrics %s failed to register", name)
			return err
		}
	} else {
		log.Errorf("Metric %s already exist", name)
		return fmt.Errorf("metric %s already exist", name)
	}
	return nil
}

func NewMetrics(name string, mtype pkg.MetricType) *Metrics {
	m := new(Metrics)
	m.Name = name
	m.Type = mtype
	return m
}

func (m *Metrics) InitializeMetric(name string, config pkg.KPIConfig, customLables []string) {
	switch config.Type {
	case pkg.MetricGuage:
		m.collector = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:        config.Name,
				Help:        config.Details,
				ConstLabels: m.Labels,
			},
			customLables,
		)

	case pkg.MetricCounter:
		m.collector = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        config.Name,
				Help:        config.Details,
				ConstLabels: m.Labels,
			},
			customLables,
		)

	case pkg.MetricSummary:
		m.collector = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:        config.Name,
				Help:        config.Details,
				ConstLabels: m.Labels,
			},
			customLables,
		)

	case pkg.MetricHistogram:
		m.collector = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        config.Name,
				Help:        config.Details,
				ConstLabels: m.Labels,
			},
			customLables,
		)
	}
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
		return fmt.Errorf("unknown metric type %s", m.Type)
	}
	return nil
}

func (m *Metrics) RegisterMetric(registry *prometheus.Registry) error {
	err := registry.Register(m.collector)
	if err != nil {
		log.Errorf("Failed to register metric %s. Err:", m.Name, err.Error())
		return err
	}
	return nil
}

func (m *Metrics) MergeLabels(static map[string]string, clabels map[string]string) {
	m.Labels = make(prometheus.Labels)
	for name, value := range static {
		m.Labels[name] = value
	}

	for name, value := range clabels {
		m.Labels[name] = value
	}
}
