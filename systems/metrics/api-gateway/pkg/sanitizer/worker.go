/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package sanitizer

import (
	"time"

	"github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/metrics/api-gateway/pkg"

	log "github.com/sirupsen/logrus"
	pmetric "github.com/ukama/ukama/systems/common/metrics"
)

const (
	NodeSubscribers = "node_subscribers"
	CountType       = "gauge"
)

var nodeMetrics = []pmetric.MetricConfig{
	{
		Name:   NodeSubscribers,
		Type:   CountType,
		Labels: map[string]string{"node_id": "", "site": "", "network": ""},
		Value:  0,
	},
}

type job struct {
	Id              int
	metricId        string
	value           float64
	task            jobFunc
	pushGatewayHost string
	nodeClient      registry.NodeClient
}

type result struct {
	Id        int
	IsSuccess bool
	Err       error
}

type metricSanitizer struct {
	nodeClient registry.NodeClient
	period     time.Duration
	stop       chan bool
}

func NewMetricSanitizer(registryHost string, period time.Duration) *metricSanitizer {
	m := &metricSanitizer{
		nodeClient: registry.NewNodeClient(registryHost),
		period:     period,
	}

	m.stop = make(chan bool)

	m.Start()

	return m
}

func (m *metricSanitizer) Start() {
	log.Infof("Starting metric sanitizer routine with period %s.", m.period)
	m.monitor()
}

func (m *metricSanitizer) Stop() {
	log.Infof("Stoping metric sanitizer routine with period %s.", m.period)

	m.stop <- true
}

func (m *metricSanitizer) sanitize() error {
	return nil
}

func (m *metricSanitizer) monitor() {
	t := time.NewTicker(m.period)

	go func() {
		for {
			select {
			case <-t.C:
				_ = m.sanitize()
			case <-m.stop:
				t.Stop()
				return
			}
		}
	}()
}

func worker(id int, jobs <-chan job, results chan<- result) {
	for j := range jobs {
		log.Info("worker", id, "started  job", j)

		err := j.task(j)

		res := result{
			Id:        j.Id,
			Err:       err,
			IsSuccess: err == nil,
		}

		log.Info("worker", id, "finished job", j)

		results <- res
	}
}

func pushNodeSubscriberMetrics(value float64, labels map[string]string, pushGatewayHost string) error {
	log.Infof("Collecting and pushing metric to push gateway host: %s", pushGatewayHost)

	return pmetric.CollectAndPushSimMetrics(pushGatewayHost, nodeMetrics,
		NodeSubscribers, float64(value), labels, pkg.SystemName)
}

type jobFunc func(job) error

func nodeSubscriberTask(j job) error {
	nodeDetails, err := j.nodeClient.Get(j.metricId)
	if err == nil {
		labels := map[string]string{
			"node_id": j.metricId,
			"network": nodeDetails.Site.NetworkId,
			"site":    nodeDetails.Site.SiteId,
		}

		err = pushNodeSubscriberMetrics(float64(j.value), labels, j.pushGatewayHost)
	}

	return err
}
