/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package collector

import (
	"fmt"
	"net/http"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	pc "github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/metrics/sanitizer/pkg"
	"google.golang.org/grpc"
)

type MetricsCollector struct {
	orgName     string
	MetricsMap  map[string]Metrics
	Config      map[string]pkg.MetricConfig
	registry    *prometheus.Registry
	grpcMetrics *grpc_prometheus.ServerMetrics
}

func NewMetricsCollector(orgName string, config []pkg.MetricConfig) *MetricsCollector {
	c := new(MetricsCollector)
	c.MetricsMap = make(map[string]Metrics)
	c.Config = make(map[string]pkg.MetricConfig, len(c.Config))
	c.grpcMetrics = grpc_prometheus.NewServerMetrics()
	c.registry = prometheus.NewRegistry()
	c.orgName = orgName
	c.registry.MustRegister(pc.NewGoCollector(), pc.NewProcessCollector(pc.ProcessCollectorOpts{}), c.grpcMetrics)

	for _, cfg := range config {
		evt := msgbus.PrepareRoute(orgName, cfg.Event)
		c.Config[evt] = cfg
	}

	return c
}

func (c *MetricsCollector) RegisterGrpcService(s *grpc.Server) {
	c.grpcMetrics.InitializeMetrics(s)
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
		log.Errorf("Event %s not expected by sanitizer service.", event)
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
