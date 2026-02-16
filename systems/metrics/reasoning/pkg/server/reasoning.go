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
	jobTag           = "reasoning-job"
	errInvalidNodeID = "Invalid node ID: %v"
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

func NewReasoningServer(msgBus mb.MsgBusServiceClient, nodeClient creg.NodeClient, config *pkg.Config, store *store.Store, sched ...scheduler.ReasoningScheduler) *ReasoningServer {
	var reasoningScheduler scheduler.ReasoningScheduler
	if len(sched) > 0 && sched[0] != nil {
		reasoningScheduler = sched[0]
	} else {
		reasoningScheduler = scheduler.NewReasoningScheduler(config.SchedulerInterval)
	}
	c := &ReasoningServer{
		store:              store,
		msgbus:             msgBus,
		config:             config,
		reasoningScheduler: reasoningScheduler,
		nodeClient:         nodeClient,
		baseRoutingKey:     msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(config.OrgName).SetService(pkg.ServiceName),
	}

	if err := c.reasoningScheduler.Start(jobTag, c.reasoningJobCallback); err != nil {
		log.Errorf("Failed to start the initial scheduler: %v", err)
	}

	return c
}

// ensureMetricKeyMap loads MetricKeyMap if nil. Returns error for gRPC handlers, nil for job (logs and continues).
func (c *ReasoningServer) ensureMetricKeyMap() error {
	if c.config.MetricKeyMap != nil {
		return nil
	}
	mkm, err := pkg.LoadMetricKeyMap(c.config)
	if err != nil {
		return err
	}
	c.config.MetricKeyMap = mkm
	return nil
}

func buildMetricQueries(metrics []pkg.Metric, nodeID string) ([]metric.MetricWithFilters, string) {
	queries := make([]metric.MetricWithFilters, 0, len(metrics))
	step := "1"
	for i, m := range metrics {
		if i == 0 && m.Step > 0 {
			step = strconv.Itoa(m.Step)
		}
		queries = append(queries, metric.MetricWithFilters{
			Metric:  m.MetricKey,
			Filters: []metric.Filter{{Key: "node_id", Value: nodeID}},
		})
	}
	return queries, step
}

func statsToProto(s algos.Stats) *pb.GetAlgoStatsForMetricResponse {
	return &pb.GetAlgoStatsForMetricResponse{
		Aggregated: &pb.AggregatedStats{
			ComputedAt:     s.ComputedAt,
			Value:          s.AggregationStats.AggregatedValue,
			Min:            s.AggregationStats.Min,
			Max:            s.AggregationStats.Max,
			P95:            s.AggregationStats.P95,
			Mean:           s.AggregationStats.Mean,
			Median:         s.AggregationStats.Median,
			SampleCount:    s.AggregationStats.SampleCount,
			Aggregation:    s.Aggregation,
			NoiseEstimate:  s.AggregationStats.NoiseEstimate,
		},
		Trend:       &pb.Trend{Type: s.Trend, Value: s.AggregationStats.AggregatedValue},
		Confidence: &pb.Confidence{Value: s.Confidence},
		Projection:  &pb.Projection{Type: s.Projection.Type, EtaSec: s.Projection.EtaSec},
		State:       s.State,
	}
}

func (c *ReasoningServer) StartScheduler(ctx context.Context, req *pb.StartSchedulerRequest) (*pb.StartSchedulerResponse, error) {
	if err := c.reasoningScheduler.Start(jobTag, c.reasoningJobCallback); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to start the scheduler: %v", err)
	}
	return &pb.StartSchedulerResponse{}, nil
}

func (c *ReasoningServer) StopScheduler(ctx context.Context, req *pb.StopSchedulerRequest) (*pb.StopSchedulerResponse, error) {
	if err := c.reasoningScheduler.Stop(jobTag); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to stop the scheduler: %v", err)
	}
	return &pb.StopSchedulerResponse{}, nil
}

func (c *ReasoningServer) reasoningJobCallback() {
	c.ReasoningJob(context.Background())
}

func (c *ReasoningServer) GetAlgoStatsForMetric(ctx context.Context, req *pb.GetAlgoStatsForMetricRequest) (*pb.GetAlgoStatsForMetricResponse, error) {
	log.Infof("Getting algo stats for metric %s on node %s", req.Metric, req.NodeId)
	nodeId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, errInvalidNodeID, err)
	}

	nType := ukama.GetNodeType(nodeId.String())
	if nType == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Could not determine node type from node ID %s", nodeId.String())
	}

	if err := c.ensureMetricKeyMap(); err != nil {
		return nil, status.Errorf(codes.Unavailable, "MetricKeyMap not loaded: %v", err)
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
	return statsToProto(stats), nil
}

func (c *ReasoningServer) GetDomains(ctx context.Context, req *pb.GetDomainsRequest) (*pb.GetDomainsResponse, error) {
	log.Infof("Getting domains for node %s", req.NodeId)
	nodeId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, errInvalidNodeID, err)
	}
	if ukama.GetNodeType(nodeId.String()) == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Could not determine node type from node ID %s", nodeId.String())
	}

	if err := c.ensureMetricKeyMap(); err != nil {
		return nil, status.Errorf(codes.Unavailable, "MetricKeyMap not loaded: %v", err)
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
		return nil, status.Errorf(codes.InvalidArgument, errInvalidNodeID, err)
	}

	nodeLog := log.WithField("node_id", req.NodeId).WithField("metric", req.Metric)
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

	now := time.Now().Unix()
	tNodeNow := now
	if e, ok := tNodeEvals[req.Metric]; ok {
		tNodeNow = e.EvaluatedAt
	}
	aNodeNow := now
	if e, ok := aNodeEvals[req.Metric]; ok {
		aNodeNow = e.EvaluatedAt
	}

	tNodeDomain := c.evaluateDomainWithEvals(nds.TNode, domain, req.Metric, tNodeEvals, rulesForMetric, tNodeNow, nodeLog)
	aNodeDomain := c.evaluateDomainWithEvals(nds.ANode, domain, req.Metric, aNodeEvals, rulesForMetric, aNodeNow, nodeLog)

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
// DomainSnapshot is not persisted, so previous is always nil (antiflap/holding only applies within a request).
func (c *ReasoningServer) evaluateDomainWithEvals(nodeID, domain, metricPattern string, evals algos.MetricEvaluationsMap, rules []algos.Rule, now int64, nodeLog *log.Entry) algos.DomainSnapshot {
	snap := algos.EvaluateDomain(domain, evals, rules, nil, now)
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

func (c *ReasoningServer) ReasoningJob(ctx context.Context) {
	jobLog := log.WithField("job", jobTag)
	jobLog.Info("Reasoning job started")

	if err := c.ensureMetricKeyMap(); err != nil {
		jobLog.WithError(err).Error("MetricKeyMap not loaded, skipping")
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

		jobLog.WithFields(log.Fields{"node_id": node.Id, "tnode": nds.TNode, "anode": nds.ANode}).Info("Processing node pair")

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

	nodeLog := jobLog.WithFields(log.Fields{"node_id": nodeID, "type": *nType})
	metricQueries, step := buildMetricQueries(metricsCfg.Metrics, nodeID)
	rp := metric.BuildPrometheusRequest(c.config.PrometheusHost, start, end, step, "", metricQueries)
	reqCtx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()
	pr, err := metric.ProcessPromRequest(reqCtx, rp)
	if err != nil {
		nodeLog.WithError(err).Error("Prometheus request failed")
		return 0, len(metricsCfg.Metrics)
	}

	for _, m := range metricsCfg.Metrics {
		prForMetric := pr.FilterResultsByMetric(m.MetricKey)
		stats, err := c.processAlgorithms(nodeID, m, prForMetric, end, nodeLog)
		if err != nil {
			nodeLog.WithError(err).WithField("metric", m.MetricKey).Error("Metric processing failed")
			failed++
		} else {
			err = c.store.PutJson(utils.GetAlgoStatsStoreKey(nodeID, m.MetricKey), stats)
			if err != nil {
				nodeLog.WithError(err).WithField("metric", m.MetricKey).Error("Failed to store stats")
				failed++
			}
			processed++
		}
	}
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

func (c *ReasoningServer) loadRules(nodeLog *log.Entry) []algos.Rule {
	paths := []string{c.config.MetricsRulesFile, "metric-rules.json", "./metric-rules.json", "../../metric-rules.json"}
	for _, p := range paths {
		rules, err := algos.LoadRulesFromJSON(p)
		if err == nil {
			nodeLog.WithField("path", p).Debug("Loaded metric rules")
			return rules
		}
		nodeLog.WithError(err).WithField("path", p).Debug("Failed to load rules from path")
	}
	nodeLog.Warn("No metric rules loaded from any path, domain evaluation will return healthy")
	return nil
}
