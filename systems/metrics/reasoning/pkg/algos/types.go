/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package algos

type Stats struct {
	Aggregation string
	AggregationStats AggregationStats
	Trend string
	Confidence float64
	State StateThresholds
}

type StatAnalysis struct {
	NewStats Stats
	PrevStats Stats
}

