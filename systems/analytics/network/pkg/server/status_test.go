/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"testing"

	"github.com/tj/assert"
)

func TestDeriveNetworkStatus(t *testing.T) {
	tests := []struct {
		name           string
		criticalAlarms int64
		openAlarms     int64
		sitesTotal     int64
		sitesDegraded  int64
		sitesOffline   int64
		want           string
	}{
		{
			name:       "all healthy",
			sitesTotal: 10,
			want:       StatusHealthy,
		},
		{
			name:           "critical alarm",
			criticalAlarms: 1,
			openAlarms:     1,
			sitesTotal:     10,
			want:           StatusCritical,
		},
		{
			name:         "more than 20 percent sites offline",
			sitesTotal:   10,
			sitesOffline: 3,
			want:         StatusCritical,
		},
		{
			name:         "exactly 20 percent sites offline is not critical",
			sitesTotal:   10,
			sitesOffline: 2,
			want:         StatusDegraded,
		},
		{
			name:       "open warning alarm",
			openAlarms: 1,
			sitesTotal: 10,
			want:       StatusDegraded,
		},
		{
			name:          "degraded site",
			sitesTotal:    10,
			sitesDegraded: 1,
			want:          StatusDegraded,
		},
		{
			name:         "single offline site",
			sitesTotal:   10,
			sitesOffline: 1,
			want:         StatusDegraded,
		},
		{
			name: "empty network is healthy",
			want: StatusHealthy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeriveNetworkStatus(tt.criticalAlarms, tt.openAlarms,
				tt.sitesTotal, tt.sitesDegraded, tt.sitesOffline)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeriveRecommendation(t *testing.T) {
	tests := []struct {
		name           string
		status         string
		offlineSeconds float64
		want           string
	}{
		{
			name:   "online is none",
			status: "online",
			want:   RecommendationNone,
		},
		{
			name:           "offline more than 30 minutes escalates",
			status:         "offline",
			offlineSeconds: 31 * 60,
			want:           RecommendationEscalate,
		},
		{
			name:           "offline less than 30 minutes restarts",
			status:         "offline",
			offlineSeconds: 10 * 60,
			want:           RecommendationRestart,
		},
		{
			name:   "needs attention restarts",
			status: "needs_attention",
			want:   RecommendationRestart,
		},
		{
			name:   "degraded restarts",
			status: "degraded",
			want:   RecommendationRestart,
		},
		{
			name:   "configuring is none",
			status: "configuring",
			want:   RecommendationNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeriveRecommendation(tt.status, tt.offlineSeconds)
			assert.Equal(t, tt.want, got)
		})
	}
}
