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
	"strconv"
	"strings"
	"time"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/algos"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/metric"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/store"
	"github.com/ukama/ukama/systems/metrics/reasoning/scheduler"
	"github.com/ukama/ukama/systems/metrics/reasoning/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func (c *ReasoningServer) GetAlgoStatsForMetric(ctx context.Context, req *pb.GetAlgoStatsForMetricRequest) (*pb.GetAlgoStatsForMetricResponse, error) {
	nodeId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid node ID: %v", err)
	}

	// TODO: Check if the metric Key is valid by checking the MetricKeyMap

	stats, err := c.loadPrevStats(utils.GetAlgoStatsStoreKey(nodeId.String(), req.Metric), log.WithField("node_id", nodeId.String()).WithField("metric", req.Metric))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Failed to load stats for metric %s on node %s: %v", req.Metric, nodeId.String(), err)
	}
	return &pb.GetAlgoStatsForMetricResponse{
		Aggregated: &pb.AggregatedStats{
			Value: stats.AggregationStats.AggregatedValue,
			Min: stats.AggregationStats.Min,
			Max: stats.AggregationStats.Max,
			P95: stats.AggregationStats.P95,
			Mean: stats.AggregationStats.Mean,
			Median: stats.AggregationStats.Median,
			SampleCount: stats.AggregationStats.SampleCount,
			Aggregation: stats.Aggregation,
			NoiseEstimate: stats.AggregationStats.NoiseEstimate,
		},
		Trend: &pb.Trend{
			Type:  stats.Trend,
			Value: stats.AggregationStats.AggregatedValue,
		},
		Confidence: &pb.Confidence{
			Value: stats.Confidence,
		},
		Projection: &pb.Projection{
			Type: stats.Projection.Type,
			EtaSec: stats.Projection.EtaSec,
		},
		State: stats.State,
	}, nil
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
	jobLog := log.WithField("job", jobTag)
	jobLog.Info("Reasoning job started")

	if c.config.MetricKeyMap == nil {
		jobLog.Error("MetricKeyMap not loaded, skipping")
		return
	}

	nodes, err := c.nodeClient.List(creg.ListNodesRequest{
		Connectivity: ukama.NodeConnectivityOnline.String(),
		State:        ukama.NodeStateConfigured.String(),
		Type:         ukama.NodeType(ukama.NODE_ID_TYPE_TOWERNODE).String(),
	})
	if err != nil {
		jobLog.WithError(err).Error("Failed to list nodes")
		return
	}
	jobLog.Infof("Processing %d node(s)", len(nodes.Nodes))

	startTime := time.Now()
	processed, failed := 0, 0

	for _, node := range nodes.Nodes {
		nds, err := utils.SortNodeIds(node.Id)
		if err != nil {
			jobLog.WithError(err).WithField("node_id", node.Id).Error("Invalid node ID, skipping")
			failed++
			continue
		}
		start, end, err := utils.GetStartEndFromStore(c.store, node.Id, c.config.PrometheusInterval)
		if err != nil {
			jobLog.WithError(err).WithField("node_id", node.Id).Error("Failed to get time window, skipping")
			failed++
			continue
		}
		for _, nodeID := range []string{nds.TNode, nds.ANode} {
			p, f := c.processNode(ctx, nodeID, start, end, jobLog)
			processed += p
			failed += f
		}
	}

	jobLog.WithFields(log.Fields{
		"processed": processed,
		"failed":   failed,
		"duration": time.Since(startTime).Round(time.Millisecond),
	}).Info("Reasoning job summary")
}

func (c *ReasoningServer) processNode(ctx context.Context, nodeID, start, end string, jobLog *log.Entry) (processed, failed int) {
	nType := ukama.GetNodeType(nodeID)
	metrics, ok := (*c.config.MetricKeyMap)[*nType]
	if !ok {
		jobLog.WithFields(log.Fields{"node_id": nodeID, "type": nType}).Debug("No metrics for node type, skipping")
		return 0, 0
	}

	nodeLog := jobLog.WithFields(log.Fields{"node_id": nodeID, "window": start + ".." + end})
	nodeLog.Debugf("Processing %d metric(s)", len(metrics.Metrics))

	for _, m := range metrics.Metrics {
		if err := c.processMetric(ctx, nodeID, m, start, end, nodeLog); err != nil {
			nodeLog.WithError(err).WithField("metric", m.Key).Error("Metric processing failed")
			failed++
		} else {
			processed++
		}
	}
	return processed, failed
}

func (c *ReasoningServer) processMetric(ctx context.Context, nodeID string, m pkg.Metric, start, end string, nodeLog *log.Entry) error {
	rp := metric.BuildPrometheusRequest(
		c.config.PrometheusHost, start, end,
		strconv.Itoa(m.Step), "",
		[]metric.MetricWithFilters{{Metric: m.Key, Filters: []metric.Filter{{Key: "node_id", Value: nodeID}}}},
	)

	pr, err := metric.ProcessPromRequest(ctx, rp)
	if err != nil {
		return err
	}

	return c.processAlgorithms(nodeID, m, pr, nodeLog)
}

func (c *ReasoningServer) processAlgorithms(nodeID string, m pkg.Metric, pr *metric.FilteredPrometheusResponse, nodeLog *log.Entry) error {
	metricLog := nodeLog.WithField("metric", m.Key)

	// Aggregate
	aggStats, err := algos.AggregateMetricAlgo(pr.Data.Result, "mean")
	if err != nil {
		return err
	}
	aggStats.RoundOfDecimalPoints(c.config.FormatDecimalPoints)

	// Load previous stats (or use empty for first run)
	storeKey := utils.GetAlgoStatsStoreKey(nodeID, m.Key)
	prevStats, err := c.loadPrevStats(storeKey, metricLog)
	if err != nil {
		return err
	}

	// Run algo pipeline
	stats := &algos.Stats{Aggregation: "mean", AggregationStats: aggStats}
	stateThresholds := algos.BuildStateThresholds(m)

	stats.Trend, err = algos.CalculateTrend(aggStats, prevStats.AggregationStats, m.TrendSensitivity)
	if err != nil {
		return err
	}

	stats.State, err = algos.CalculateState(aggStats.AggregatedValue, stateThresholds, m.StateDirection)
	if err != nil {
		return err
	}

	expectedSamples := 0
	if m.Step > 0 {
		expectedSamples = c.config.PrometheusInterval / m.Step
	}
	stats.Confidence, err = algos.CalculateConfidence(pr.Data.Result, aggStats, prevStats.AggregationStats, stats.State, expectedSamples)
	if err != nil {
		return err
	}
	stats.Confidence = utils.RoundToDecimalPoints(stats.Confidence, c.config.FormatDecimalPoints)

	stats.Projection = algos.ProjectCrossingTime(
		aggStats.AggregatedValue, prevStats.AggregationStats.AggregatedValue,
		float64(c.config.PrometheusInterval), stateThresholds, m.StateDirection,
	)
	if stats.Projection.Type != "" {
		stats.Projection.EtaSec = utils.RoundToDecimalPoints(stats.Projection.EtaSec, c.config.FormatDecimalPoints)
	}

	// Persist
	if err := c.store.PutJson(storeKey, stats); err != nil {
		metricLog.WithError(err).Error("Failed to persist algo stats")
		return err
	}

	log.Infof("Metics Data: %+v", pr.Data.Result)
	c.logMetricResult(nodeID, m.Key, stats, metricLog)
	return nil
}

func (c *ReasoningServer) loadPrevStats(storeKey string, metricLog *log.Entry) (algos.Stats, error) {
	bytes, err := c.store.GetJson(storeKey)
	if err == nil {
		stats, err := algos.UnmarshalStatsFromJSON(bytes)
		if err != nil {
			return algos.Stats{}, err
		}
		return stats, nil
	}
	if strings.Contains(err.Error(), "not found") {
		metricLog.Info("No previous stats found, using empty (first run or new metric)")
		return algos.EmptyPrevStats(), nil
	}
	return algos.Stats{}, err
}

func (c *ReasoningServer) logMetricResult(nodeID, metricKey string, stats *algos.Stats, metricLog *log.Entry) {
	d := utils.MetricResultLogData{
		Value:            stats.AggregationStats.AggregatedValue,
		State:            stats.State,
		Trend:            stats.Trend,
		Confidence:       stats.Confidence,
		ProjectionType:   stats.Projection.Type,
		ProjectionEtaSec: stats.Projection.EtaSec,
	}
	metricLog.WithFields(utils.MetricResultLogFields(d)).Infof("%s/%s: %s (%s)", nodeID, metricKey, stats.State, stats.Trend)
}

