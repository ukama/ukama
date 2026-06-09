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
	"time"

	"github.com/stretchr/testify/assert"
)

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestDiagnose(t *testing.T) {
	now := time.Date(2026, 6, 3, 12, 0, 0, 0, time.UTC)
	recentlySeen := timePtr(now.Add(-10 * time.Minute))

	healthy := diagnosisInput{
		Now:               now,
		SimStatus:         "active",
		PackageStatus:     "active",
		HasSiteHealth:     true,
		SiteUptimePercent: 99.9,
		UsageLast24hMb:    120,
		LastSeenAt:        recentlySeen,
	}

	tests := []struct {
		name             string
		mutate           func(in diagnosisInput) diagnosisInput
		wantIssue        string
		wantAction       string
		wantEscalation   bool
		wantSignalStates map[string]string
	}{
		{
			name:           "all healthy",
			mutate:         func(in diagnosisInput) diagnosisInput { return in },
			wantIssue:      "No issue detected",
			wantAction:     "No action needed",
			wantEscalation: false,
			wantSignalStates: map[string]string{
				signalKeySimState:     signalStateOk,
				signalKeyPackageState: signalStateOk,
				signalKeySiteHealth:   signalStateOk,
				signalKeyUsage:        signalStateOk,
				signalKeyLastSeen:     signalStateOk,
			},
		},
		{
			name: "faulty sim is critical and escalates",
			mutate: func(in diagnosisInput) diagnosisInput {
				in.SimStatus = "faulty"
				return in
			},
			wantIssue:      "Faulty SIM",
			wantAction:     "Replace SIM",
			wantEscalation: true,
			wantSignalStates: map[string]string{
				signalKeySimState: signalStateCritical,
			},
		},
		{
			name: "suspended sim is critical",
			mutate: func(in diagnosisInput) diagnosisInput {
				in.SimStatus = "suspended"
				return in
			},
			wantIssue:      "SIM is suspended",
			wantAction:     "Review suspension and reactivate SIM",
			wantEscalation: true,
			wantSignalStates: map[string]string{
				signalKeySimState: signalStateCritical,
			},
		},
		{
			name: "expired package is warning",
			mutate: func(in diagnosisInput) diagnosisInput {
				in.PackageStatus = "expired"
				return in
			},
			wantIssue:      "Package expired",
			wantAction:     "Renew package",
			wantEscalation: false,
			wantSignalStates: map[string]string{
				signalKeyPackageState: signalStateWarning,
			},
		},
		{
			name: "site uptime below 90 is critical and dominates package warning",
			mutate: func(in diagnosisInput) diagnosisInput {
				in.PackageStatus = "expired"
				in.SiteUptimePercent = 80
				return in
			},
			wantIssue:      "Site is down or unstable",
			wantAction:     "Escalate to network team",
			wantEscalation: true,
			wantSignalStates: map[string]string{
				signalKeyPackageState: signalStateWarning,
				signalKeySiteHealth:   signalStateCritical,
			},
		},
		{
			name: "site uptime below 99 is warning",
			mutate: func(in diagnosisInput) diagnosisInput {
				in.SiteUptimePercent = 95
				return in
			},
			wantIssue:      "Site is degraded",
			wantAction:     "Escalate to network team",
			wantEscalation: false,
			wantSignalStates: map[string]string{
				signalKeySiteHealth: signalStateWarning,
			},
		},
		{
			name: "zero usage is warning with check device action",
			mutate: func(in diagnosisInput) diagnosisInput {
				in.UsageLast24hMb = 0
				return in
			},
			wantIssue:      "No recent data usage",
			wantAction:     "Check device",
			wantEscalation: false,
			wantSignalStates: map[string]string{
				signalKeyUsage: signalStateWarning,
			},
		},
		{
			name: "stale last seen is warning",
			mutate: func(in diagnosisInput) diagnosisInput {
				in.LastSeenAt = timePtr(now.Add(-48 * time.Hour))
				return in
			},
			wantIssue:      "Customer not seen recently",
			wantAction:     "Check device",
			wantEscalation: false,
			wantSignalStates: map[string]string{
				signalKeyLastSeen: signalStateWarning,
			},
		},
		{
			name: "missing site health is unknown but not an issue",
			mutate: func(in diagnosisInput) diagnosisInput {
				in.HasSiteHealth = false
				in.SiteUptimePercent = 0
				return in
			},
			wantIssue:      "No issue detected",
			wantAction:     "No action needed",
			wantEscalation: false,
			wantSignalStates: map[string]string{
				signalKeySiteHealth: signalStateUnknown,
			},
		},
		{
			name: "sim critical wins over earlier warnings",
			mutate: func(in diagnosisInput) diagnosisInput {
				in.SimStatus = "faulty"
				in.PackageStatus = "expired"
				in.UsageLast24hMb = 0
				return in
			},
			wantIssue:      "Faulty SIM",
			wantAction:     "Replace SIM",
			wantEscalation: true,
			wantSignalStates: map[string]string{
				signalKeySimState:     signalStateCritical,
				signalKeyPackageState: signalStateWarning,
				signalKeyUsage:        signalStateWarning,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := diagnose(tt.mutate(healthy))

			assert.Equal(t, tt.wantIssue, res.LikelyIssue)
			assert.Equal(t, tt.wantAction, res.RecommendedAction)
			assert.Equal(t, tt.wantEscalation, res.EscalationNeeded)
			assert.Len(t, res.Signals, 5)

			got := make(map[string]string)
			for _, s := range res.Signals {
				got[s.Key] = s.State
			}

			for key, state := range tt.wantSignalStates {
				assert.Equal(t, state, got[key], "signal %s", key)
			}
		})
	}
}
