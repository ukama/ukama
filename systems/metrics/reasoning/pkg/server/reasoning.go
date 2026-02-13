/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/algos"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/metric"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/store"
	"github.com/ukama/ukama/systems/metrics/reasoning/scheduler"
	"github.com/ukama/ukama/systems/metrics/reasoning/utils"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	ukama "github.com/ukama/ukama/systems/common/ukama"
	pb "github.com/ukama/ukama/systems/metrics/reasoning/pb/gen"
)

const (
	jobTag               = "reasoning-job"
)

type ReasoningServer struct {
	pb.UnimplementedReasoningServiceServer
	msgbus             mb.MsgBusServiceClient
	baseRoutingKey     msgbus.RoutingKeyBuilder
	reasoningScheduler scheduler.ReasoningScheduler
	config             *pkg.Config
	nodeClient         creg.NodeClient
	store    		   *store.Store
}

func NewReasoningServer(msgBus mb.MsgBusServiceClient, nodeClient creg.NodeClient, config *pkg.Config, store *store.Store) *ReasoningServer {
	scheduler := scheduler.NewReasoningScheduler(config.SchedulerInterval)
	c := &ReasoningServer{
		store:              store,
		msgbus:             msgBus,
		config:             config,
		reasoningScheduler: scheduler,
		nodeClient:         nodeClient,
		baseRoutingKey:  msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(config.OrgName).SetService(pkg.ServiceName),
	}

	// Start the scheduler
	if err := c.reasoningScheduler.Start(jobTag, func() { c.ReasoningJob(context.Background()) }); err != nil {
		log.Errorf("Failed to start the initial scheduler: %v", err)
	}

	return c
}

func (c *ReasoningServer) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	return &pb.GetStatsResponse{}, nil
}

func (c *ReasoningServer) GetDomains(ctx context.Context, req *pb.GetDomainsRequest) (*pb.GetDomainsResponse, error) {
	return &pb.GetDomainsResponse{}, nil
}

func (c *ReasoningServer) StartScheduler(ctx context.Context, req *pb.StartSchedulerRequest) (*pb.StartSchedulerResponse, error) {
	log.Info("Starting scheduler")

	return &pb.StartSchedulerResponse{}, nil
}

func (c *ReasoningServer) StopScheduler(ctx context.Context, req *pb.StopSchedulerRequest) (*pb.StopSchedulerResponse, error) {
	log.Info("Stopping scheduler")

	return &pb.StopSchedulerResponse{}, nil
}

func (c *ReasoningServer) ReasoningJob(ctx context.Context) {
	log.Info("Reasoning job started")
	if c.config.MetricKeyMap == nil {
		log.Error("MetricKeyMap not loaded")
		return
	}

	nodes, err := c.nodeClient.List(creg.ListNodesRequest{
		Connectivity: ukama.NodeConnectivityOnline.String(),
		State:        ukama.NodeStateConfigured.String(),
		Type:         ukama.NodeType(ukama.NODE_ID_TYPE_TOWERNODE).String(),
	})
	if err != nil {
		log.Errorf("Failed to get nodes: %v", err)
		return
	}
	log.Infof("Node registry nodes: %v", nodes.Nodes)

	for _, node := range nodes.Nodes {
		nds, err := utils.SortNodeIds(node.Id)
		if err != nil {
			log.Errorf("Failed to sort node IDs for %s: %v", node.Id, err)
			continue
		}
		start, end, err := utils.GetStartEndFromStore(c.store, node.Id, c.config.PrometheusInterval)
		if err != nil {
			log.Errorf("Failed to get start/end for node %s: %v", node.Id, err)
			continue
		}
		for _, nodeID := range []string{nds.TNode, nds.ANode} {
			c.processNode(ctx, nodeID, start, end)
		}
	}
}

func (c *ReasoningServer) processNode(ctx context.Context, nodeID string, start, end string) {
	nType := ukama.GetNodeType(nodeID)
	metrics, ok := (*c.config.MetricKeyMap)[*nType]
	if !ok {
		log.Debugf("No metrics for node type %v, skipping %s", nType, nodeID)
		return
	}
	log.Debugf("Processing %d metrics for node %s", len(metrics.Metrics), nodeID)

	for _, m := range metrics.Metrics {
		if err := c.processMetric(ctx, nodeID, m, start, end); err != nil {
			log.Errorf("Metric %s for node %s: %v", m.Key, nodeID, err)
		}
	}
}

func (c *ReasoningServer) processMetric(ctx context.Context, nodeID string, m pkg.Metric, start, end string) error {
	rp := metric.BuildPrometheusRequest(
		c.config.PrometheusHost,
		start, end,
		strconv.Itoa(m.Step),
		"",
		[]metric.MetricWithFilters{{Metric: m.Key, Filters: []metric.Filter{{Key: "node_id", Value: nodeID}}}},
	)
	log.Debugf("Prometheus request: %s - %s{%s}", rp.Url, m.Key, nodeID)

	pr, err := metric.ProcessPromRequest(ctx, rp)
	if err != nil {
		return err
	}

	return c.processAlgorithms(nodeID, m, pr)
}

func (c *ReasoningServer) processAlgorithms(nodeID string, m pkg.Metric, pr *metric.FilteredPrometheusResponse) error {
	statAnalysis := algos.StatAnalysis{}
	err := errors.New("error processing algorithms")
	
	statAnalysis.NewStats.AggregationStats, err = algos.AggregateMetricAlgo(pr.Data.Result, "mean")
	if err != nil {
		return err
	}
	statAnalysis.NewStats.AggregationStats.RoundOfDecimalPoints(c.config.FormatDecimalPoints)

	storeKey := utils.GetAlgoStatsStoreKey(nodeID, m.Key)
	prevAggStatsBytes, err := c.store.GetJson(storeKey)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			// First run or new metric: use empty prev stats so algorithms still run
			statAnalysis.PrevStats = algos.EmptyPrevStats()
		} else {
			return err
		}
	} else {
		statAnalysis.PrevStats, err = algos.UnmarshalStatsFromJSON(prevAggStatsBytes)
		if err != nil {
			return err
		}
	}

	statAnalysis.NewStats.Trend, err = algos.CalculateTrend(statAnalysis.NewStats.AggregationStats, statAnalysis.PrevStats.AggregationStats, m.TrendSensitivity)
	if err != nil {
		return err
	}

	stateThresholds := algos.BuildStateThresholds(m)
	statAnalysis.NewStats.State, err = algos.CalculateState(
		statAnalysis.NewStats.AggregationStats.AggregatedValue,
		stateThresholds,
		m.StateDirection,
	)
	if err != nil {
		return err
	}

	expectedSamples := 0
	if m.Step > 0 {
		expectedSamples = c.config.PrometheusInterval / m.Step
	}
	statAnalysis.NewStats.Confidence, err = algos.CalculateConfidence(
		pr.Data.Result,
		statAnalysis.NewStats.AggregationStats,
		statAnalysis.PrevStats.AggregationStats,
		statAnalysis.NewStats.State,
		expectedSamples,
	)
	if err != nil {
		return err
	}
	statAnalysis.NewStats.Confidence = utils.RoundToDecimalPoints(statAnalysis.NewStats.Confidence, c.config.FormatDecimalPoints)

	statAnalysis.NewStats.Projection = algos.ProjectCrossingTime(
		statAnalysis.NewStats.AggregationStats.AggregatedValue,
		statAnalysis.PrevStats.AggregationStats.AggregatedValue,
		float64(c.config.PrometheusInterval),
		stateThresholds,
		m.StateDirection,
	)
	if statAnalysis.NewStats.Projection.Type != "" {
		statAnalysis.NewStats.Projection.EtaSec = utils.RoundToDecimalPoints(statAnalysis.NewStats.Projection.EtaSec, c.config.FormatDecimalPoints)
	} else {
		log.Debugf("No projection for %s/%s: valueNow=%.2f valuePrev=%.2f warning=%.2f direction=%s (needs two windows + slope toward threshold)",
			nodeID, m.Key,
			statAnalysis.NewStats.AggregationStats.AggregatedValue,
			statAnalysis.PrevStats.AggregationStats.AggregatedValue,
			stateThresholds.Warning,
			m.StateDirection,
		)
	}

	if err := c.store.PutJson(storeKey, statAnalysis.NewStats); err != nil {
		log.Errorf("Failed to persist algo stats for %s/%s: %v", nodeID, m.Key, err)
		return err
	}

	log.Infof("Stat analysis: %+v", statAnalysis)
	return nil
}

