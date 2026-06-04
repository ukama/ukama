/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

// Pure status/recommendation derivations, kept free of I/O for testability.

const (
	StatusHealthy  = "healthy"
	StatusDegraded = "degraded"
	StatusCritical = "critical"

	RecommendationRestart  = "restart"
	RecommendationEscalate = "escalate"
	RecommendationNone     = "none"
)

// offlineEscalateSeconds: a resource offline for more than 30 minutes should
// be escalated rather than restarted.
const offlineEscalateSeconds = 30 * 60

// DeriveNetworkStatus derives the overall network status:
//   - critical: any open critical alarm, or more than 20% of sites offline
//   - degraded: any open alarm, or any site degraded/offline
//   - healthy: otherwise
func DeriveNetworkStatus(criticalAlarms, openAlarms int64, sitesTotal, sitesDegraded, sitesOffline int64) string {
	if criticalAlarms > 0 {
		return StatusCritical
	}

	if sitesTotal > 0 && float64(sitesOffline)/float64(sitesTotal) > 0.2 {
		return StatusCritical
	}

	if openAlarms > 0 || sitesDegraded > 0 || sitesOffline > 0 {
		return StatusDegraded
	}

	return StatusHealthy
}

// DeriveRecommendation derives a support recommendation for a site or node:
//   - escalate: offline for more than 30 minutes
//   - restart: status is needs_attention, degraded or offline (recent)
//   - none: otherwise
func DeriveRecommendation(status string, offlineDurationSeconds float64) string {
	if status == "offline" && offlineDurationSeconds > offlineEscalateSeconds {
		return RecommendationEscalate
	}

	if status == "needs_attention" || status == "degraded" || status == "offline" {
		return RecommendationRestart
	}

	return RecommendationNone
}
