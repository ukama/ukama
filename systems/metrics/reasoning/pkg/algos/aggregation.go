/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package algos

import (
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/metric"
)

type AggregationStats struct {
	Min         float64
	Max         float64
	P95         float64
	Mean        float64
	Median      float64
	SampleCount float64
	Aggregation string
}

func AggregateMetricResults(results []metric.FilteredPrometheusResult) ([]metric.FilteredPrometheusResult, error) {
	return results, nil
}