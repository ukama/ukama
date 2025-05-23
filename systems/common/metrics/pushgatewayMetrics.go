/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	log "github.com/sirupsen/logrus"
)

type MetricConfig struct {
	Name    string
	Type    string
	Labels  map[string]string
	Details string
	Buckets []float64
	Value   float64
}

type MetricType string

const (
	MetricUnknown   MetricType = "unavailable"
	MetricGauge     MetricType = "gauge"
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
		return MetricGauge
	case "counter":
		return MetricCounter
	case "histogram":
		return MetricHistogram
	case "summary":
		return MetricSummary
	default:
		return MetricUnknown
	}
}

var metricCollectors = make(map[string][]*Metrics)

func NewMetrics(name string, mtype string) *Metrics {
	m := new(Metrics)
	m.Name = name
	m.Type = MetricTypeFromString(mtype)
	return m
}

func (m *Metrics) InitializeMetric(name string, config MetricConfig, customLabels []string) error {
	switch MetricTypeFromString(config.Type) {
	case MetricGauge:
		m.collector = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:        config.Name,
				Help:        config.Details,
				ConstLabels: m.Labels,
			},
			customLabels,
		)

	case MetricCounter:
		m.collector = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        config.Name,
				Help:        config.Details,
				ConstLabels: m.Labels,
			},
			customLabels,
		)

	case MetricSummary:
		m.collector = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:        config.Name,
				Help:        config.Details,
				ConstLabels: m.Labels,
			},
			customLabels,
		)

	case MetricHistogram:
		m.collector = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        config.Name,
				Help:        config.Details,
				ConstLabels: m.Labels,
				Buckets:     config.Buckets,
			},
			customLabels,
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

func PushMetrics(pushMetricHost string, metrics []MetricConfig, metricJobName string) {
	log.Infof("Pushing metric job: %s", metricJobName)

	labelDimensions := make([]string, 0, len(metrics[0].Labels))
	for key := range metrics[0].Labels {
		labelDimensions = append(labelDimensions, key)
	}

	for _, metric := range metrics {
		if _, ok := metricCollectors[metric.Name]; !ok {
			// Metric does not exist, create a new one
			newMetric := NewMetrics(metric.Name, metric.Type)
			if err := newMetric.InitializeMetric(metric.Name, MetricConfig{
				Name:   metric.Name,
				Type:   metric.Type,
				Labels: metric.Labels,
			}, labelDimensions); err != nil {
				log.Errorf("Failed to initialize metric %s: %v", metric.Name, err)
				continue
			}
			metricCollectors[metric.Name] = []*Metrics{newMetric}
		}

		for _, m := range metricCollectors[metric.Name] {
			if err := m.SetMetric(metric.Value, metric.Labels); err != nil {
				log.Errorf("Failed to set metric %s value: %v", metric.Name, err)
				continue
			}
		}
	}

	pusher := push.New(pushMetricHost, metricJobName)
	for _, metrics := range metricCollectors {
		for _, m := range metrics {
			pusher.Collector(m.collector)
		}
	}
	if err := pusher.Push(); err != nil {
		log.Errorf("Could not push metrics to Pushgateway: %s", err.Error())
	}
}

func CollectAndPushSimMetrics(pushGateway string, configMetrics []MetricConfig, selectedMetric string, Value float64, Labels map[string]string, systemName string) error {
	log.Infof("Collecting and pushing metric %q on behalf of system %q", selectedMetric, systemName)

	var selectedMetrics []MetricConfig
	var foundSelectedMetric bool

	for i, metric := range configMetrics {
		if metric.Name == selectedMetric {
			metric.Value = Value
			for k, v := range Labels {
				metric.Labels[k] = v
			}
			selectedMetrics = append(selectedMetrics, metric)
			configMetrics[i] = metric
			foundSelectedMetric = true
			break
		}
	}

	if !foundSelectedMetric {
		log.Errorf("metric %q not found", selectedMetric)

		return fmt.Errorf("metric %q not found", selectedMetric)
	}

	PushMetrics(pushGateway, selectedMetrics, systemName)

	return nil
}
