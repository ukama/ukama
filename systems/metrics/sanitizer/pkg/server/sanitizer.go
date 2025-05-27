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

	// "github.com/klauspost/compress/snappy"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/metrics/sanitizer/pkg"

	log "github.com/sirupsen/logrus"
	pmetric "github.com/ukama/ukama/systems/common/metrics"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/metrics/sanitizer/pb/gen"
)

const (
	name = "__name__"
	env  = "env"
	job  = "job"
	// mainLabel = "nodeId"
	mainLabel = "network"
)

type NodeMetaData struct {
	NodeId    string
	NetworkId string
	SiteId    string
}

type NodeMetricMetaData struct {
	MainLabelValue   string
	AdditionalLabels map[string]string
	Value            float64
}

type SanitizerServer struct {
	pb.UnimplementedSanitizerServiceServer
	baseRoutingKey  msgbus.RoutingKeyBuilder
	registryHost    string
	pushGatewayHost string
	NodeCache       map[string]NodeMetaData
	NodeMetricCache map[string]float64
	org             string
	orgName         string
	msgbus          mb.MsgBusServiceClient
}

func NewSanitizerServer(registryHost, pushGatewayHost, orgName string, org string,
	msgBus mb.MsgBusServiceClient) (*SanitizerServer, error) {
	s := SanitizerServer{
		registryHost:    registryHost,
		pushGatewayHost: pushGatewayHost,
		orgName:         orgName,
		org:             org,
		msgbus:          msgBus,
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

	s.NodeMetricCache = map[string]float64{}

	return &s, nil
}

func (s *SanitizerServer) Sanitize(ctx context.Context, req *pb.SanitizeRequest) (*pb.SanitizeResponse, error) {
	log.Info("Getting a sanitize request")

	var metricsPayload prompb.WriteRequest

	metricsToPush := []NodeMetricMetaData{}

	data, err := snappy.Decode(nil, req.Data)
	if err != nil {
		log.Errorf("Fail to decode remote_write data. Error: %v", err)

		return nil, fmt.Errorf("fail to decode remote_write data. Error: %w", err)
	}

	err = metricsPayload.Unmarshal(data)
	if err != nil {
		log.Errorf("Fail to unmarshal remote_write data. Error: %v", err)

		return nil, fmt.Errorf("fail to unmarshal remote_write data. Error: %w", err)
	}

	// data, err = json.Marshal(metricsPayload)
	// if err != nil {
	// return nil, err
	// }

	// log.Infof("Raw body: %s", string(data))

	for _, ts := range metricsPayload.Timeseries {
		metric := NodeMetricMetaData{
			AdditionalLabels: make(map[string]string)}

		if len(ts.Samples) > 0 {
			metric.Value = ts.Samples[0].Value

			for _, label := range ts.Labels {
				if label.Name == name || label.Name == env || label.Name == job {
					continue
				}

				if label.Name == mainLabel {
					metric.MainLabelValue = label.Value
					continue
				}

				metric.AdditionalLabels[label.Name] = label.Value
			}

			if metric.MainLabelValue == "" {
				log.Warnf("main label %q not found in timeseries data, moving on to next metric...",
					mainLabel)

				continue
			}

			value, ok := s.NodeMetricCache[metric.MainLabelValue]
			if !ok || value != metric.Value {
				s.NodeMetricCache[metric.MainLabelValue] = metric.Value

				cachedNode, ok := s.NodeCache[metric.MainLabelValue]
				if !ok {
					log.Warnf("metadata not found in cache for nodeId: %s, skipping...",
						metric.MainLabelValue)

					continue
				}

				metric.AdditionalLabels["network"] = cachedNode.NetworkId
				metric.AdditionalLabels["site"] = cachedNode.SiteId
				metricsToPush = append(metricsToPush, metric)
			}
		}
	}

	for _, m := range metricsToPush {
		pushUpdatedNodeMetrics(m.Value, m.AdditionalLabels, s.pushGatewayHost)
	}

	return nil, nil
}

func (s *SanitizerServer) syncNodeCache() error {
	log.Infof("Fetching list of nodes with metadata.")

	nCache := map[string]NodeMetaData{}

	nodeClient := registry.NewNodeClient(s.registryHost)
	resp, err := nodeClient.GetAll()
	if err != nil {
		log.Errorf("Fail to get list of nodes with metadata: Error: %v", err)

		return fmt.Errorf("fail to get list of nodes with metadata: Error: %w", err)
	}

	for _, n := range resp.Nodes {
		if n.Site.NodeId != "" {
			nCache[n.Site.NodeId] = NodeMetaData{
				NodeId:    n.Site.NodeId,
				NetworkId: n.Site.NetworkId,
				SiteId:    n.Site.SiteId,
			}
		}
	}

	s.NodeCache = nCache

	return nil
}

func pushUpdatedNodeMetrics(value float64, labels map[string]string, pushGatewayHost string) {
	log.Infof("Collecting and pushing node active subscribers metric to push gateway host: %s", pushGatewayHost)

	err := pmetric.CollectAndPushSimMetrics(pushGatewayHost, pkg.NodeActiveSubscribersMetric,
		pkg.NodeActiveSubscribers, float64(value), labels, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while pushing node active subscribers metric to push gateway %s", err.Error())
	}
}
