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
)

// EmptyPrevStats returns Stats for first run when no previous data exists.
// Uses NaN for AggregatedValue so trend/projection/confidence treat it as "no prior data".
func EmptyPrevStats() Stats {
	return Stats{
		AggregationStats: AggregationStats{AggregatedValue: math.NaN()},
	}
}

type Stats struct {
	Aggregation      string `json:"aggregation"`
	AggregationStats AggregationStats `json:"aggregation_stats"`
	Trend            string `json:"trend"`
	Confidence       float64 `json:"confidence"`
	State            string `json:"state"`
	Projection       ProjectionStats `json:"projection"`
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
