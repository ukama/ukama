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
	"testing"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/systems/common/ukama"
	pb "github.com/ukama/ukama/systems/metrics/reasoning/pb/gen"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/algos"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/store"
	"github.com/ukama/ukama/systems/metrics/reasoning/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mockScheduler implements scheduler.ReasoningScheduler for testing.
type mockScheduler struct {
	startErr error
	stopErr  error
}

func (m *mockScheduler) SetNewJob(tag string, taskFunc any, params ...any) (*gocron.Job, error) {
	return nil, nil
}

func (m *mockScheduler) Start(tag string, taskFunc any, params ...any) error {
	return m.startErr
}

func (m *mockScheduler) Stop(tag string) error {
	return m.stopErr
}

func testMetricKeyMap() pkg.MetricKeyMap {
	return pkg.MetricKeyMap{
		ukama.NODE_TYPE_TOWERNODE: pkg.Metrics{
			Metrics: []pkg.Metric{
				{Step: 1, Key: "cpu", MetricKey: "com_soc_cpu_usage", Category: "health", TrendSensitivity: 1, StateDirection: "higher_is_worse"},
				{Step: 1, Key: "memory", MetricKey: "com_memory_ddr_used", Category: "health", TrendSensitivity: 1, StateDirection: "higher_is_worse"},
			},
		},
		ukama.NODE_TYPE_AMPNODE: pkg.Metrics{
			Metrics: []pkg.Metric{
				{Step: 1, Key: "cpu", MetricKey: "ctl_soc_cpu_usage", Category: "health", TrendSensitivity: 1, StateDirection: "higher_is_worse"},
				{Step: 1, Key: "memory", MetricKey: "ctl_memory_ddr_used", Category: "health", TrendSensitivity: 1, StateDirection: "higher_is_worse"},
			},
		},
	}
}

func testConfig() *pkg.Config {
	mkm := testMetricKeyMap()
	return &pkg.Config{
		OrgName:            "ukama",
		PrometheusInterval:  60,
		SchedulerInterval:  24 * time.Hour,
		FormatDecimalPoints: 2,
		MetricsRulesFile:   "metric-rules.json",
		MetricKeyMap:       &mkm,
	}
}

func validTowerNodeID() string {
	return ukama.NewVirtualTowerNodeId().String()
}

func sampleStats() algos.Stats {
	return algos.Stats{
		Aggregation: "mean",
		AggregationStats: algos.AggregationStats{
			AggregatedValue: 25.5,
			Min:            10,
			Max:            40,
			P95:            38,
			Mean:           25.5,
			Median:         24,
			SampleCount:    60,
			NoiseEstimate:  2,
		},
		Trend:       "stable",
		Confidence:  0.85,
		State:       "healthy",
		Projection:  algos.ProjectionStats{},
		ComputedAt:  time.Now().Unix(),
	}
}

func TestStartScheduler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock := &mockScheduler{startErr: nil, stopErr: nil}
		cfg := testConfig()
		str := store.NewInMemoryStore()
		srv := NewReasoningServer(nil, nil, cfg, str, mock)

		ctx := context.Background()
		resp, err := srv.StartScheduler(ctx, &pb.StartSchedulerRequest{})
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
	t.Run("Error", func(t *testing.T) {
		mock := &mockScheduler{startErr: errors.New("scheduler start failed"), stopErr: nil}
		cfg := testConfig()
		str := store.NewInMemoryStore()
		srv := NewReasoningServer(nil, nil, cfg, str, mock)

		ctx := context.Background()
		resp, err := srv.StartScheduler(ctx, &pb.StartSchedulerRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to start the scheduler")
	})
}

func TestStopScheduler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock := &mockScheduler{startErr: nil, stopErr: nil}
		cfg := testConfig()
		str := store.NewInMemoryStore()
		srv := NewReasoningServer(nil, nil, cfg, str, mock)

		ctx := context.Background()
		resp, err := srv.StopScheduler(ctx, &pb.StopSchedulerRequest{})
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
	t.Run("Error", func(t *testing.T) {
		mock := &mockScheduler{startErr: nil, stopErr: errors.New("scheduler stop failed")}
		cfg := testConfig()
		str := store.NewInMemoryStore()
		srv := NewReasoningServer(nil, nil, cfg, str, mock)

		ctx := context.Background()
		resp, err := srv.StopScheduler(ctx, &pb.StopSchedulerRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to stop the scheduler")
	})
}

func TestGetAlgoStatsForMetric(t *testing.T) {
	t.Run("InvalidNodeID", func(t *testing.T) {
		mock := &mockScheduler{}
		cfg := testConfig()
		str := store.NewInMemoryStore()
		srv := NewReasoningServer(nil, nil, cfg, str, mock)

		ctx := context.Background()
		resp, err := srv.GetAlgoStatsForMetric(ctx, &pb.GetAlgoStatsForMetricRequest{
			NodeId: "invalid-node-id",
			Metric: "cpu",
		})
		require.Error(t, err)
		require.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid node ID")
	})
	t.Run("InvalidMetricKey", func(t *testing.T) {
		mock := &mockScheduler{}
		cfg := testConfig()
		str := store.NewInMemoryStore()
		srv := NewReasoningServer(nil, nil, cfg, str, mock)

		nodeID := validTowerNodeID()
		ctx := context.Background()
		resp, err := srv.GetAlgoStatsForMetric(ctx, &pb.GetAlgoStatsForMetricRequest{
			NodeId: nodeID,
			Metric: "nonexistent_metric",
		})
		require.Error(t, err)
		require.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "not valid")
	})
	t.Run("EmptyStoreReturnsSuccess", func(t *testing.T) {
		mock := &mockScheduler{}
		cfg := testConfig()
		str := store.NewInMemoryStore()
		srv := NewReasoningServer(nil, nil, cfg, str, mock)

		nodeID := validTowerNodeID()
		ctx := context.Background()
		resp, err := srv.GetAlgoStatsForMetric(ctx, &pb.GetAlgoStatsForMetricRequest{
			NodeId: nodeID,
			Metric: "cpu",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Aggregated)
	})
	t.Run("Success", func(t *testing.T) {
		mock := &mockScheduler{}
		cfg := testConfig()
		str := store.NewInMemoryStore()
		nodeID := validTowerNodeID()
		storeKey := utils.GetAlgoStatsStoreKey(nodeID, "com_soc_cpu_usage")
		stats := sampleStats()
		err := str.PutJson(storeKey, &stats)
		require.NoError(t, err)

		srv := NewReasoningServer(nil, nil, cfg, str, mock)

		ctx := context.Background()
		resp, err := srv.GetAlgoStatsForMetric(ctx, &pb.GetAlgoStatsForMetricRequest{
			NodeId: nodeID,
			Metric: "cpu",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Aggregated)
		assert.Equal(t, stats.AggregationStats.AggregatedValue, resp.Aggregated.Value)
		assert.Equal(t, stats.AggregationStats.Mean, resp.Aggregated.Mean)
		assert.Equal(t, stats.State, resp.State)
		assert.Equal(t, stats.Trend, resp.Trend.Type)
		assert.Equal(t, stats.Confidence, resp.Confidence.Value)
	})
	t.Run("AnodeNodeID", func(t *testing.T) {
		mock := &mockScheduler{}
		cfg := testConfig()
		str := store.NewInMemoryStore()
		nodeID := ukama.NewVirtualAmplifierNodeId().String()
		storeKey := utils.GetAlgoStatsStoreKey(nodeID, "ctl_soc_cpu_usage")
		stats := sampleStats()
		err := str.PutJson(storeKey, &stats)
		require.NoError(t, err)

		srv := NewReasoningServer(nil, nil, cfg, str, mock)

		ctx := context.Background()
		resp, err := srv.GetAlgoStatsForMetric(ctx, &pb.GetAlgoStatsForMetricRequest{
			NodeId: nodeID,
			Metric: "cpu",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Aggregated)
		assert.Equal(t, stats.AggregationStats.AggregatedValue, resp.Aggregated.Value)
	})
	t.Run("NoMetricKeyMap", func(t *testing.T) {
		mock := &mockScheduler{}
		cfg := testConfig()
		cfg.MetricKeyMap = nil
		cfg.MetricsKeyMapFile = "/nonexistent/path"
		str := store.NewInMemoryStore()
		nodeID := validTowerNodeID()
		srv := NewReasoningServer(nil, nil, cfg, str, mock)

		ctx := context.Background()
		resp, err := srv.GetAlgoStatsForMetric(ctx, &pb.GetAlgoStatsForMetricRequest{
			NodeId: nodeID,
			Metric: "cpu",
		})
		require.Error(t, err)
		require.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Unavailable, st.Code())
		assert.Contains(t, st.Message(), "MetricKeyMap not loaded")
	})
}

func TestGetDomains(t *testing.T) {
	t.Run("InvalidNodeID", func(t *testing.T) {
		mock := &mockScheduler{}
		cfg := testConfig()
		str := store.NewInMemoryStore()
		srv := NewReasoningServer(nil, nil, cfg, str, mock)

		ctx := context.Background()
		resp, err := srv.GetDomains(ctx, &pb.GetDomainsRequest{
			NodeId: "invalid",
			Metric: "cpu",
		})
		require.Error(t, err)
		require.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})
	t.Run("InvalidMetric", func(t *testing.T) {
		mock := &mockScheduler{}
		cfg := testConfig()
		str := store.NewInMemoryStore()
		nodeID := validTowerNodeID()
		srv := NewReasoningServer(nil, nil, cfg, str, mock)

		ctx := context.Background()
		resp, err := srv.GetDomains(ctx, &pb.GetDomainsRequest{
			NodeId: nodeID,
			Metric: "invalid_metric",
		})
		require.Error(t, err)
		require.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "not valid")
	})
	t.Run("Success", func(t *testing.T) {
		mock := &mockScheduler{}
		cfg := testConfig()
		str := store.NewInMemoryStore()
		nodeID := validTowerNodeID()
		nds, err := utils.SortNodeIds(nodeID)
		require.NoError(t, err)

		stats := sampleStats()
		storeKeyT := utils.GetAlgoStatsStoreKey(nds.TNode, "com_soc_cpu_usage")
		storeKeyA := utils.GetAlgoStatsStoreKey(nds.ANode, "ctl_soc_cpu_usage")
		err = str.PutJson(storeKeyT, &stats)
		require.NoError(t, err)
		err = str.PutJson(storeKeyA, &stats)
		require.NoError(t, err)

		srv := NewReasoningServer(nil, nil, cfg, str, mock)

		ctx := context.Background()
		resp, err := srv.GetDomains(ctx, &pb.GetDomainsRequest{
			NodeId: nodeID,
			Metric: "cpu",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Domain)
		assert.NotEmpty(t, resp.Domain.Severity)
		assert.NotEmpty(t, resp.Domain.RuleId)
	})
	t.Run("NoMetricKeyMap", func(t *testing.T) {
		mock := &mockScheduler{}
		cfg := testConfig()
		cfg.MetricKeyMap = nil
		cfg.MetricsKeyMapFile = "/nonexistent/path"
		str := store.NewInMemoryStore()
		nodeID := validTowerNodeID()
		srv := NewReasoningServer(nil, nil, cfg, str, mock)

		ctx := context.Background()
		resp, err := srv.GetDomains(ctx, &pb.GetDomainsRequest{
			NodeId: nodeID,
			Metric: "cpu",
		})
		require.Error(t, err)
		require.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Unavailable, st.Code())
		assert.Contains(t, st.Message(), "MetricKeyMap not loaded")
	})
}
