/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func RecordSuccessfulRequestMetric() {
	go func() {
		opsSuccessProcessed.Inc()
	}()
}

func RecordFailedRequestMetric() {
	go func() {
		opsFailedProcessed.Inc()
	}()
}

func RecordRetiredRequestMetric() {
	go func() {
		opsRetired.Inc()
	}()
}

func RecordSupersededRequestMetric() {
	go func() {
		opsSuperseded.Inc()
	}()
}

var (
	opsRetired = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "node_feeder_requests_dropped_total",
		Help:        "Requests dropped instead of retried",
		ConstLabels: map[string]string{"reason": "retired"},
	})
	opsSuperseded = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "node_feeder_requests_dropped_total",
		Help:        "Requests dropped instead of retried",
		ConstLabels: map[string]string{"reason": "superseded"},
	})
)

var (
	opsSuccessProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "node_feeder_requests_total",
		Help:        "The total number requests",
		ConstLabels: map[string]string{"status": "succeeded"},
	})
)

var (
	opsFailedProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "node_feeder_requests_total",
		Help:        "The total number requests",
		ConstLabels: map[string]string{"status": "failed"},
	})
)
