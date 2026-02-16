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
	"encoding/json"
	"strconv"
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
	jobTag = "reasoning-job"
)

type metricWithStats struct {
	metric pkg.Metric
	stats  *algos.Stats
}

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
	log.Infof("Getting algo stats for metric %s on node %s", req.Metric, req.NodeId)
	nodeId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid node ID: %v", err)
	}

	nType := ukama.GetNodeType(nodeId.String())
	if nType == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Could not determine node type from node ID %s", nodeId.String())
	}

	if c.config.MetricKeyMap == nil {
		metricKeyMap, err := pkg.LoadMetricKeyMap(c.config)
		if err != nil {
			return nil, status.Errorf(codes.Unavailable, "MetricKeyMap not loaded: %v", err)
		}
		c.config.MetricKeyMap = metricKeyMap
	}

	metricsCfg, ok := (*c.config.MetricKeyMap)[*nType]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "No metrics configured for node type %s", *nType)
	}

	metricKey, err := utils.ValidateMetricKey(req.Metric, metricsCfg, *nType)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Metric key %q is not valid for node type %s: %v", req.Metric, *nType, err)
	}

	stats, err := algos.LoadStats(c.store, utils.GetAlgoStatsStoreKey(nodeId.String(), metricKey), log.WithField("node_id", nodeId.String()).WithField("metric", metricKey))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Failed to load stats for metric %s on node %s: %v", req.Metric, nodeId.String(), err)
	}
	return &pb.GetAlgoStatsForMetricResponse{
		Aggregated: &pb.AggregatedStats{
			ComputedAt: stats.ComputedAt,
			Value:      stats.AggregationStats.AggregatedValue,
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
	log.Infof("Getting domains for node %s", req.NodeId)
	nodeId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid node ID: %v", err)
	}
	if ukama.GetNodeType(nodeId.String()) == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Could not determine node type from node ID %s", nodeId.String())
	}

	if c.config.MetricKeyMap == nil {
		metricKeyMap, err := pkg.LoadMetricKeyMap(c.config)
		if err != nil {
			return nil, status.Errorf(codes.Unavailable, "MetricKeyMap not loaded: %v", err)
		}
		c.config.MetricKeyMap = metricKeyMap
	}

	tNodeMetrics, ok := (*c.config.MetricKeyMap)[ukama.NODE_TYPE_TOWERNODE]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "No metrics configured for node type tnode")
	}
	aNodeMetrics, ok := (*c.config.MetricKeyMap)[ukama.NODE_TYPE_AMPNODE]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "No metrics configured for node type anode")
	}

	if _, err := utils.ValidateMetricKey(req.Metric, tNodeMetrics, ukama.NODE_TYPE_TOWERNODE); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Metric key %q is not valid: %v", req.Metric, err)
	}

	nds, err := utils.SortNodeIds(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid node ID: %v", err)
	}

	now := time.Now().Unix()
	nodeLog := log.WithField("node_id", req.NodeId).WithField("metric", req.Metric)

	// Load stats for ALL metric types for both TNode and ANode, then build MetricEvaluationsMap for domain evaluation
	tNodeEvals := c.buildMetricEvaluationsForNode(nds.TNode, tNodeMetrics, nodeLog)
	aNodeEvals := c.buildMetricEvaluationsForNode(nds.ANode, aNodeMetrics, nodeLog)

	rules := c.loadRules(nodeLog)
	rulesForMetric := algos.RulesForMetric(rules, req.Metric)
	if len(rulesForMetric) == 0 {
		nodeLog.Debug("No rules for metric, returning healthy")
	}

	domain := "health"
	if len(tNodeMetrics.Metrics) > 0 && tNodeMetrics.Metrics[0].Category != "" {
		domain = tNodeMetrics.Metrics[0].Category
	}

	tNodeDomain := c.evaluateDomainWithEvals(nds.TNode, domain, req.Metric, tNodeEvals, rulesForMetric, now, nodeLog)
	aNodeDomain := c.evaluateDomainWithEvals(nds.ANode, domain, req.Metric, aNodeEvals, rulesForMetric, now, nodeLog)

	if algos.SeverityRank(tNodeDomain.Severity) >= algos.SeverityRank(aNodeDomain.Severity) {
		return &pb.GetDomainsResponse{Domain: domainSnapshotToProto(&tNodeDomain)}, nil
	}
	return &pb.GetDomainsResponse{Domain: domainSnapshotToProto(&aNodeDomain)}, nil
}

// buildMetricEvaluationsForNode loads stats for all metrics and builds MetricEvaluationsMap (pattern key -> MetricEvaluation).
func (c *ReasoningServer) buildMetricEvaluationsForNode(nodeID string, metricsCfg pkg.Metrics, nodeLog *log.Entry) algos.MetricEvaluationsMap {
	evals := make(algos.MetricEvaluationsMap)
	for _, m := range metricsCfg.Metrics {
		stats, err := algos.LoadStats(c.store, utils.GetAlgoStatsStoreKey(nodeID, m.MetricKey), nodeLog.WithField("metric", m.MetricKey))
		if err != nil {
			nodeLog.WithError(err).WithField("metric", m.MetricKey).Debug("No stats for metric, skipping from evals")
			continue
		}
		evals[m.Key] = algos.MetricEvaluationFromStats(m.Key, stats)
	}
	return evals
}

// evaluateDomainWithEvals runs domain evaluation given MetricEvaluationsMap and filtered rules.
func (c *ReasoningServer) evaluateDomainWithEvals(nodeID, domain, metricPattern string, evals algos.MetricEvaluationsMap, rules []algos.Rule, now int64, nodeLog *log.Entry) algos.DomainSnapshot {
	var previous *algos.DomainSnapshot
	if domainKey := utils.GetDomainStoreKey(nodeID, metricPattern); domainKey != "" {
		if prev, err := c.loadDomainSnapshot(domainKey, nodeLog); err == nil && prev.RuleID != "" {
			previous = prev
		}
	}
	snap := algos.EvaluateDomain(domain, evals, rules, previous, now)
	nodeLog.WithFields(log.Fields{"node_id": nodeID, "domain": domain, "rule_id": snap.RuleID, "severity": snap.Severity}).Debug("Domain evaluation")
	return snap
}

func domainSnapshotToProto(s *algos.DomainSnapshot) *pb.Domain {
	if s == nil {
		return nil
	}
	return &pb.Domain{
		RuleId:         s.RuleID,
		Severity:       s.Severity,
		Headline:       s.Headline,
		RootCause:      s.RootCause,
		ServiceImpact:  s.ServiceImpact,
		RuleConfidence: s.RuleConfidence,
		EvaluatedAt:    s.EvaluatedAt,
		ComputedAt:     s.EvaluatedAt,
	}
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

	// Retry loading MetricKeyMap if nil (e.g. file wasn't available at startup)
	if c.config.MetricKeyMap == nil {
		metricKeyMap, err := pkg.LoadMetricKeyMap(c.config)
		if err != nil {
			jobLog.WithError(err).Error("MetricKeyMap not loaded, skipping")
			return
		}
		c.config.MetricKeyMap = metricKeyMap
		jobLog.Info("MetricKeyMap loaded successfully")
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
		jobLog.WithFields(log.Fields{"node_id": node.Id, "tnode": nds.TNode, "anode": nds.ANode}).Info("Processing tower+amp node pair")

		jobLog.WithField("node_id", node.Id).Info("Getting time window from store")
		start, end, err := utils.GetStartEndFromStore(c.store, node.Id, c.config.PrometheusInterval)
		if err != nil {
			jobLog.WithError(err).WithField("node_id", node.Id).Error("Failed to get time window, skipping")
			failed++
			continue
		}
		jobLog.WithFields(log.Fields{"node_id": node.Id, "start": start, "end": end}).Info("Got time window, processing metrics")

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
	if nType == nil {
		jobLog.WithField("node_id", nodeID).Warn("Could not determine node type from node ID, skipping")
		return 0, 0
	}
	metricsCfg, ok := (*c.config.MetricKeyMap)[*nType]
	if !ok {
		jobLog.WithFields(log.Fields{"node_id": nodeID, "type": *nType}).Info("No metrics for node type in MetricKeyMap, skipping")
		return 0, 0
	}

	nodeLog := jobLog.WithFields(log.Fields{"node_id": nodeID, "type": *nType, "window": start + ".." + end})
	nodeLog.Infof("Processing %d metric(s) for node", len(metricsCfg.Metrics))

	// Batch all metrics into a single Prometheus request (per node_types config)
	metricQueries := make([]metric.MetricWithFilters, 0, len(metricsCfg.Metrics))
	step := "1"
	for i, m := range metricsCfg.Metrics {
		if i == 0 && m.Step > 0 {
			step = strconv.Itoa(m.Step)
		}
		metricQueries = append(metricQueries, metric.MetricWithFilters{
			Metric:  m.MetricKey,
			Filters: []metric.Filter{{Key: "node_id", Value: nodeID}},
		})
	}

	rp := metric.BuildPrometheusRequest(c.config.PrometheusHost, start, end, step, "", metricQueries)
	nodeLog.WithFields(log.Fields{"metrics": len(metricQueries), "prometheus_host": c.config.PrometheusHost}).Debug("Fetching metrics from Prometheus")
	reqCtx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()
	pr, err := metric.ProcessPromRequest(reqCtx, rp)
	if err != nil {
		nodeLog.WithError(err).Error("Prometheus request failed")
		return 0, len(metricsCfg.Metrics)
	}

	// Process each metric and collect stats (no persist yet)
	// collected := make([]metricWithStats, 0, len(metricsCfg.Metrics))
	for _, m := range metricsCfg.Metrics {
		prForMetric := pr.FilterResultsByMetric(m.MetricKey)
		s, err := c.processAlgorithms(nodeID, m, prForMetric, end, nodeLog)
		if err != nil {
			nodeLog.WithError(err).WithField("metric", m.MetricKey).Error("Metric processing failed")
			failed++
		} else {
			c.store.PutJson(utils.GetAlgoStatsStoreKey(nodeID, m.MetricKey), s)
			processed++
		}
		// collected = append(collected, metricWithStats{metric: m, stats: s})
	}

	// Domain evaluation and persist: store algo stats and domain separately
	// if len(collected) > 0 {
	// 	domainSnapshot := c.evaluateDomainForNode(nodeID, collected, end, nodeLog)
	// 	for _, ms := range collected {
			
	// 		domainKey := utils.GetDomainStoreKey(nodeID, ms.metric.Key)
	// 		if err := c.store.PutJson(domainKey, domainSnapshot); err != nil {
	// 			nodeLog.WithError(err).WithField("metric", ms.metric.Key).Error("Failed to persist domain")
	// 		}
	// 	}
	// }
	return processed, failed
}

func (c *ReasoningServer) processAlgorithms(nodeID string, m pkg.Metric, pr *metric.FilteredPrometheusResponse, end string, nodeLog *log.Entry) (*algos.Stats, error) {
	metricLog := nodeLog.WithField("metric", m.MetricKey)

	aggStats, err := algos.AggregateMetricAlgo(pr.Data.Result, "mean")
	if err != nil {
		return nil, err
	}
	aggStats.RoundOfDecimalPoints(c.config.FormatDecimalPoints)

	storeKey := utils.GetAlgoStatsStoreKey(nodeID, m.MetricKey)
	prevStats, err := algos.LoadStats(c.store, storeKey, metricLog)
	if err != nil {
		return nil, err
	}

	stats := &algos.Stats{Aggregation: "mean", AggregationStats: aggStats}
	stateThresholds := algos.BuildStateThresholds(m)

	stats.Trend, err = algos.CalculateTrend(aggStats, prevStats.AggregationStats, m.TrendSensitivity)
	if err != nil {
		return nil, err
	}

	stats.State, err = algos.CalculateState(aggStats.AggregatedValue, stateThresholds, m.StateDirection)
	if err != nil {
		return nil, err
	}

	expectedSamples := 0
	if m.Step > 0 {
		expectedSamples = c.config.PrometheusInterval / m.Step
	}
	stats.Confidence, err = algos.CalculateConfidence(pr.Data.Result, aggStats, prevStats.AggregationStats, stats.State, expectedSamples)
	if err != nil {
		return nil, err
	}
	stats.Confidence = utils.RoundToDecimalPoints(stats.Confidence, c.config.FormatDecimalPoints)

	stats.Projection = algos.ProjectCrossingTime(
		aggStats.AggregatedValue, prevStats.AggregationStats.AggregatedValue,
		float64(c.config.PrometheusInterval), stateThresholds, m.StateDirection,
	)
	if stats.Projection.Type != "" {
		stats.Projection.EtaSec = utils.RoundToDecimalPoints(stats.Projection.EtaSec, c.config.FormatDecimalPoints)
	}

	endUnix, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		endUnix = time.Now().Unix()
	}
	stats.ComputedAt = endUnix

	return stats, nil
}

func (c *ReasoningServer) evaluateDomainForNode(nodeID string, collected []metricWithStats, end string, nodeLog *log.Entry) algos.DomainSnapshot {
	rules := c.loadRules(nodeLog)
	if len(rules) == 0 {
		return algos.DomainSnapshot{}
	}

	evals := make(algos.MetricEvaluationsMap)
	for _, ms := range collected {
		evals[ms.metric.Key] = algos.MetricEvaluation{
			MetricID:    ms.metric.Key,
			State:       ms.stats.State,
			Trend:       ms.stats.Trend,
			Conclusion:  algos.CombineStateAndTrend(ms.stats.State, ms.stats.Trend),
			Confidence:  ms.stats.Confidence,
			EvaluatedAt: ms.stats.ComputedAt,
		}
	}

	endUnix, _ := strconv.ParseInt(end, 10, 64)
	if endUnix == 0 {
		endUnix = time.Now().Unix()
	}

	var previous *algos.DomainSnapshot
	if len(collected) > 0 {
		domainKey := utils.GetDomainStoreKey(nodeID, collected[0].metric.Key)
		if prev, err := c.loadDomainSnapshot(domainKey, nodeLog); err == nil && prev.RuleID != "" {
			previous = prev
		}
	}

	domain := "health"
	if len(collected) > 0 && collected[0].metric.Category != "" {
		domain = collected[0].metric.Category
	}
	snap := algos.EvaluateDomain(domain, evals, rules, previous, endUnix)
	nodeLog.WithFields(log.Fields{"node_id": nodeID, "domain": domain, "rule_id": snap.RuleID, "severity": snap.Severity}).Debug("Domain evaluation")
	return snap
}

func (c *ReasoningServer) loadDomainSnapshot(storeKey string, nodeLog *log.Entry) (*algos.DomainSnapshot, error) {
	bytes, err := c.store.GetJson(storeKey)
	if err != nil {
		return nil, err
	}
	var snap algos.DomainSnapshot
	if err := json.Unmarshal(bytes, &snap); err != nil {
		return nil, err
	}
	return &snap, nil
}

func (c *ReasoningServer) loadRules(nodeLog *log.Entry) []algos.Rule {
	paths := []string{c.config.MetricsRulesFile, "metric-rules.json", "./metric-rules.json"}
	for _, p := range paths {
		rules, err := algos.LoadRulesFromJSON(p)
		if err == nil {
			return rules
		}
	}
	nodeLog.Debug("No metric rules loaded, domain evaluation skipped")
	return nil
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

