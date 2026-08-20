/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/prometheus/prometheus/prompb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/metrics/sanitizer/pkg"

	snappy "github.com/klauspost/compress/s2"
	log "github.com/sirupsen/logrus"
	pmetric "github.com/ukama/ukama/systems/common/metrics"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/metrics/sanitizer/pb/gen"
)

const (
	name = "__name__"

	// Node scrape jobs are not consistent on how they name the node label:
	// the k8s mesh-nodes job stamps `node_id` while some paths stamp `nodeid`.
	// Both are accepted as the node identifier on the way in, and the
	// republished metrics always carry `node_id`.
	nodeLabel    = pkg.NodeIdLabel
	altNodeLabel = "nodeid"

	networkLabel = pkg.NetworkLabel
	siteLabel    = pkg.SiteIdLabel
)

type NodeMetaData struct {
	NodeId    string
	NetworkId string
	SiteId    string
}

type NodeMetricMetaData struct {
	// MetricName is the sanitized name the sample is republished under.
	MetricName string
	NodeId     string
	Labels     map[string]string
	Value      float64
}

type SanitizerServer struct {
	pb.UnimplementedSanitizerServiceServer
	baseRoutingKey  msgbus.RoutingKeyBuilder
	registryHost    string
	pushGatewayHost string
	nodeCache       map[string]NodeMetaData
	nodeMetricCache map[string]float64
	org             string
	orgName         string
	msgbus          mb.MsgBusServiceClient
	m               *sync.RWMutex
	pushMutex       *sync.Mutex
}

func NewSanitizerServer(registryHost, pushGatewayHost, orgName string, org string,
	msgBus mb.MsgBusServiceClient) (*SanitizerServer, error) {
	s := SanitizerServer{
		registryHost:    registryHost,
		pushGatewayHost: pushGatewayHost,
		nodeMetricCache: map[string]float64{},
		orgName:         orgName,
		org:             org,
		msgbus:          msgBus,
		m:               &sync.RWMutex{},
		pushMutex:       &sync.Mutex{},
	}

	if msgBus != nil {
		s.baseRoutingKey = msgbus.NewRoutingKeyBuilder().SetCloudSource().
			SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName)
	}

	err := s.syncNodeCache()
	if err != nil {
		log.Errorf("error while initializing new sanitizer server: %v", err)

		return nil, fmt.Errorf("error while initializing new sanitizer server: %w", err)
	}

	return &s, nil
}

func (s *SanitizerServer) Sanitize(ctx context.Context, req *pb.SanitizeRequest) (*pb.SanitizeResponse, error) {
	log.Info("Getting a sanitize request")

	var metricsPayload prompb.WriteRequest

	metricsToPush := []NodeMetricMetaData{}

	data, err := snappy.Decode(nil, req.Data)
	if err != nil {
		log.Errorf("Failed to decode remote_write data. Error: %v", err)

		return nil, fmt.Errorf("failed to decode remote_write data. Error: %w", err)
	}

	err = metricsPayload.Unmarshal(data)
	if err != nil {
		log.Errorf("Failed to unmarshal remote_write data. Error: %v", err)

		return nil, fmt.Errorf("failed to unmarshal remote_write data. Error: %w", err)
	}

	for _, ts := range metricsPayload.Timeseries {
		if len(ts.Samples) == 0 {
			continue
		}

		rawName, nodeId := metricAndNodeId(ts.Labels)

		sanitized, ok := pkg.SanitizedMetrics[rawName]
		if !ok {
			log.Debugf("metric %q is not sanitized by this service, moving on to next metric...",
				rawName)

			continue
		}

		if nodeId == "" {
			log.Warnf("node label (%q or %q) not found on metric %q, moving on to next metric...",
				nodeLabel, altNodeLabel, rawName)

			continue
		}

		// A remote_write batch may carry several samples for the same series:
		// the most recent one is what's worth republishing.
		value := latestSampleValue(ts.Samples)

		log.Infof("processing metric %q for node %s with sample value: %v",
			rawName, nodeId, value)

		if sanitized.SkipUnchanged {
			cached, ok := s.getNodeMetricFromCache(sanitized.Name, nodeId)
			if ok && cached == value {
				log.Infof("No new value to cache for metric %q on node %s: %f, skipping ...",
					sanitized.Name, nodeId, cached)

				continue
			}
		}

		cachedNode, ok := s.getNodeFromCache(nodeId)
		if !ok {
			log.Warnf("metadata not found in cache for nodeId: %s, we'll be skipping...", nodeId)
			log.Warn("make sure all physical nodes are correctly registered under registry, nodes")

			continue
		}

		if sanitized.SkipUnchanged {
			s.updateNodeMetricCache(sanitized.Name, nodeId, value)
		}

		// Only the declared label set is forwarded: common/metrics pins a
		// metric's label dimensions on the first push, so letting arbitrary
		// scrape labels (instance, pod, ...) through would break every
		// subsequent push of that same metric.
		metricsToPush = append(metricsToPush, NodeMetricMetaData{
			MetricName: sanitized.Name,
			NodeId:     nodeId,
			Value:      value,
			Labels: map[string]string{
				nodeLabel:    nodeId,
				siteLabel:    cachedNode.SiteId,
				networkLabel: cachedNode.NetworkId,
			},
		})
	}

	for _, m := range metricsToPush {
		s.pushSanitizedNodeMetric(m)
	}

	return &pb.SanitizeResponse{}, nil
}

// latestSampleValue returns the value of the most recent sample of a series.
// Samples reach us in timestamp order, but the order is not guaranteed by the
// remote write protocol, so the timestamps are compared rather than trusted.
func latestSampleValue(samples []prompb.Sample) float64 {
	latest := samples[0]

	for _, sample := range samples[1:] {
		if sample.Timestamp > latest.Timestamp {
			latest = sample
		}
	}

	return latest.Value
}

// metricAndNodeId extracts the metric name and the node identifier out of a
// timeseries label set. Every other label is dropped: the sanitizer rebuilds
// the label set of the republished metric from the node cache.
func metricAndNodeId(labels []prompb.Label) (metricName string, nodeId string) {
	for _, label := range labels {
		switch label.Name {
		case name:
			metricName = label.Value
		case nodeLabel, altNodeLabel:
			if nodeId == "" {
				nodeId = label.Value
			}
		}
	}

	return metricName, nodeId
}

func (s *SanitizerServer) updateNodeCache(n map[string]NodeMetaData) {
	s.m.Lock()
	defer s.m.Unlock()

	s.nodeCache = n
}

func (s *SanitizerServer) getNodeFromCache(nodeId string) (NodeMetaData, bool) {
	s.m.RLock()
	defer s.m.RUnlock()

	node, ok := s.nodeCache[nodeId]

	return node, ok
}

// metricCacheKey scopes the last seen value to a metric/node pair: without the
// metric name, the several metrics sanitized for a same node would overwrite
// each other's cached value.
func metricCacheKey(metricName, nodeId string) string {
	return metricName + "/" + nodeId
}

func (s *SanitizerServer) updateNodeMetricCache(metricName, nodeId string, value float64) {
	s.m.Lock()
	defer s.m.Unlock()

	s.nodeMetricCache[metricCacheKey(metricName, nodeId)] = value
}

func (s *SanitizerServer) getNodeMetricFromCache(metricName, nodeId string) (float64, bool) {
	s.m.RLock()
	defer s.m.RUnlock()

	value, ok := s.nodeMetricCache[metricCacheKey(metricName, nodeId)]

	return value, ok
}

func (s *SanitizerServer) syncNodeCache() error {
	log.Infof("Fetching list of nodes with metadata.")

	nCache := map[string]NodeMetaData{}

	nodeClient := registry.NewNodeClient(s.registryHost)
	resp, err := nodeClient.GetAll()
	if err != nil {
		log.Errorf("Failed to get list of nodes with metadata: Error: %v", err)

		return fmt.Errorf("failed to get list of nodes with metadata: Error: %w", err)
	}

	log.Infof("Found %d node(s) to cache", len(resp.Nodes))

	for _, n := range resp.Nodes {
		if n.Site.SiteId != "" {
			nCache[n.Id] = NodeMetaData{
				NodeId:    n.Id,
				NetworkId: n.Site.NetworkId,
				SiteId:    n.Site.SiteId,
			}
		}
	}

	s.updateNodeCache(nCache)
	log.Infof("Cached %d node(s)", len(nCache))

	return nil
}

func (s *SanitizerServer) pushSanitizedNodeMetric(m NodeMetricMetaData) {
	if math.IsNaN(m.Value) {
		// Prometheus turns a series that stopped reporting into a staleness
		// marker, which travels over remote_write as a NaN sample. It is
		// forwarded as is: the sanitizer only appends labels, it does not
		// decide whether a node is up.
		log.Infof("Pushing metric %q for node %s with no value (NaN) to push gateway host: %s",
			m.MetricName, m.NodeId, s.pushGatewayHost)
	} else {
		log.Infof("Pushing metric %q for node %s to push gateway host: %s",
			m.MetricName, m.NodeId, s.pushGatewayHost)
	}

	// common/metrics keeps its collectors in a package level map and merges the
	// caller labels into the shared metric config, neither of which is
	// goroutine safe, so the pushes issued by this server are serialised.
	s.pushMutex.Lock()
	defer s.pushMutex.Unlock()

	err := pmetric.CollectAndPushSystemMetrics(s.pushGatewayHost, pkg.NodeMetrics,
		m.MetricName, m.Value, m.Labels, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while pushing metric %q for node %s to push gateway: %v",
			m.MetricName, m.NodeId, err)
	}
}
