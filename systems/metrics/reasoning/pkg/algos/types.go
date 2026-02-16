/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package algos

import (
	"encoding/json"
	"math"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/store"
)

// EmptyPrevStats returns Stats for first run when no previous data exists.
// Uses NaN for AggregatedValue so trend/projection/confidence treat it as "no prior data".
func EmptyPrevStats() Stats {
	return Stats{
		AggregationStats: AggregationStats{AggregatedValue: math.NaN()},
	}
}

type MetricEvaluation struct {
	MetricID    string          `json:"metric_id"`
	EvaluatedAt int64           `json:"evaluated_at"`
	AggNow      float64         `json:"agg_now"`
	StatsNow    AggregationStats `json:"stats_now"`
	State       string          `json:"state"`
	Trend       string          `json:"trend"`
	Conclusion  string          `json:"conclusion"`
	Confidence  float64         `json:"confidence"`
	Projection  ProjectionStats `json:"projection"`
}

type Stats struct {
	Aggregation      string          `json:"aggregation"`
	AggregationStats AggregationStats `json:"aggregation_stats"`
	Trend            string          `json:"trend"`
	Confidence       float64         `json:"confidence"`
	State            string          `json:"state"`
	Projection       ProjectionStats `json:"projection"`
	ComputedAt       int64           `json:"computed_at"`
}

type StatAnalysis struct {
	NewStats Stats `json:"new_stats"`
	PrevStats Stats `json:"prev_stats"`
}

func (s *StatAnalysis) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NewStats Stats `json:"new_stats"`
		PrevStats Stats `json:"prev_stats"`
	}{
		NewStats: s.NewStats,
		PrevStats: s.PrevStats,
	})
}

func (s *StatAnalysis) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *Stats) MarshalStatsToJSON() ([]byte, error) {
	return json.Marshal(s)
}

func UnmarshalStatsFromJSON(data []byte) (Stats, error) {
	var stats Stats
	err := json.Unmarshal(data, &stats)
	if err != nil {
		return Stats{}, err
	}
	return stats, nil
}

func LoadStats(store *store.Store, storeKey string, metricLog *log.Entry) (Stats, error) {
	bytes, err := store.GetJson(storeKey)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			metricLog.Info("No previous stats found, using empty (first run or new metric)")
			return EmptyPrevStats(), nil
		}
		return Stats{}, err
	}
	stats, err := UnmarshalStatsFromJSON(bytes)
	if err != nil {
		// Corrupted or empty stored data: treat as first run so algorithms still execute
		metricLog.WithError(err).Info("Invalid previous stats, using empty (first run or corrupted data)")
		return EmptyPrevStats(), nil
	}
	return stats, nil
}

func MetricEvaluationFromStats(metricKey string, stats Stats) MetricEvaluation {
	return MetricEvaluation{
		MetricID:    metricKey,
		State:       stats.State,
		Trend:       stats.Trend,
		Conclusion:  CombineStateAndTrend(stats.State, stats.Trend),
		Confidence:  stats.Confidence,
		EvaluatedAt: stats.ComputedAt,
	}
}