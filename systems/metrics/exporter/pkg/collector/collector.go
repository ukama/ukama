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

type MetricsCollector struct {
	MetricsMap map[string]Metrics
	Config     map[string]pkg.MetricConfig
	registry   *prometheus.Registry
}

func NewMetricsCollector(config []pkg.MetricConfig) *MetricsCollector {
	c := new(MetricsCollector)
	c.MetricsMap = make(map[string]Metrics)
	c.Config = make(map[string]pkg.MetricConfig, len(c.Config))
	c.registry = prometheus.NewRegistry()
	c.registry.MustRegister(pc.NewGoCollector(), pc.NewProcessCollector(pc.ProcessCollectorOpts{}))

	for _, cfg := range config {
		c.Config[cfg.Event] = cfg
	}

	return c
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

func (c *MetricsCollector) GetConfigForEvent(event string) (*pkg.MetricConfig, error) {

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
