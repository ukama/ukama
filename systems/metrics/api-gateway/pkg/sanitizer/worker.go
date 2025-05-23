/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package sanitizer

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
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

type jobFunc func(Job) error

type Job struct {
	Id              int
	metricId        string
	value           float64
	task            jobFunc
	pushGatewayHost string
	nodeClient      registry.NodeClient
}

type Result struct {
	Id        int
	IsSuccess bool
	Err       error
}

type metricSanitizer struct {
	m          *pkg.Metrics
	nodeClient registry.NodeClient
	period     time.Duration
	stop       chan bool
}

func NewMetricSanitizer(metricClient *pkg.Metrics, registryHost string, period time.Duration) *metricSanitizer {
	m := &metricSanitizer{
		m:          metricClient,
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
	// for each metric to sanitize: spin a goroutine that;
	//		fetch all metrics values (node_subscribers) metrics to sanitize
	w := bytes.NewBuffer([]byte{})
	resp := map[string]any{}
	jobs := make(chan Job, 1)
	results := make(chan Result, 1)

	go worker(1, jobs, results)

	err := m.requestMetric(w, "sites", pkg.NewFilter())
	if err != nil {
		log.Errorf("Failure to request metric %s. Error: %v", "sites", err)

		return err
	}

	//		unmarshell the response out of the io.writer
	err = m.unmarshalMetricFromWriter(w, resp)
	if err != nil {
		log.Errorf("Failure to unmarshal metric response %v. Error: %v", resp, err)

		return err
	}

	//		for each metric record send out the job tothe worker pool channel
	//		should the the jobs channel be local or global?
	err = m.sendJob(resp, jobs)
	if err != nil {
		log.Errorf("Failure to send job from metric response %v. Error: %v", resp, err)

		return err
	}

	// log results
	res := <-results
	log.Infof("Task result: %v", res)

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

func (m *metricSanitizer) requestMetric(writer io.ReadWriter, metric string, filter *pkg.Filter) error {
	_, err := m.m.GetMetric(strings.ToLower(metric), filter, writer, true)

	return err
}

func (m *metricSanitizer) unmarshalMetricFromWriter(reader io.Reader, resp map[string]any) error {
	err := json.NewDecoder(reader).Decode(&resp)
	if err != nil {
		log.Errorf("Failed to unmarshal metric response. Error: %v", err)

		return err
	}

	return nil
}

func (m *metricSanitizer) sendJob(metricPayload map[string]any, jobs chan<- Job) error {
	log.Printf("sending Job for payload: %v", metricPayload)

	if v, ok := metricPayload["status"]; ok {
		if v, ok := v.(string); ok {
			if v == "success" {
				j := Job{
					Id:   1,
					task: nodeSubscriberTask,
				}

				jobs <- j
				close(jobs)
			}
		}
	}

	return nil
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
	for j := range jobs {
		log.Info("worker", id, "started  job", j)

		err := j.task(j)

		res := Result{
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

func nodeSubscriberTask(j Job) error {
	log.Infof("Running task function for job %v", j)

	// nodeDetails, err := j.nodeClient.Get(j.metricId)
	// if err == nil {
	// labels := map[string]string{
	// "node_id": j.metricId,
	// "network": nodeDetails.Site.NetworkId,
	// "site":    nodeDetails.Site.SiteId,
	// }

	// err = pushNodeSubscriberMetrics(float64(j.value), labels, j.pushGatewayHost)
	// }

	// return err

	return nil
}
