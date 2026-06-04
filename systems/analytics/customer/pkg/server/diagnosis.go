/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"fmt"
	"time"
)

const (
	signalStateOk       = "ok"
	signalStateWarning  = "warning"
	signalStateCritical = "critical"
	signalStateUnknown  = "unknown"

	signalKeySimState     = "sim_state"
	signalKeyPackageState = "package_state"
	signalKeySiteHealth   = "site_health"
	signalKeyUsage        = "usage"
	signalKeyLastSeen     = "last_seen"
)

// diagnosisInput is the pure-data input for support diagnosis. It is
// deliberately decoupled from db/pb types so the diagnosis logic can be
// unit-tested without any infrastructure.
type diagnosisInput struct {
	Now time.Time

	SimStatus     string /* available/assigned/active/suspended/faulty */
	PackageStatus string /* active/expired/... */

	HasSiteHealth     bool
	SiteUptimePercent float64

	UsageLast24hMb float64
	LastSeenAt     *time.Time
}

type diagnosisSignal struct {
	Key    string
	State  string
	Detail string
}

type diagnosisResult struct {
	Signals           []diagnosisSignal
	LikelyIssue       string
	RecommendedAction string
	EscalationNeeded  bool
}

// diagnose is a PURE function deriving support signals, the likely issue and
// the recommended action from a customer's current state.
func diagnose(in diagnosisInput) diagnosisResult {
	signals := make([]diagnosisSignal, 0, 5)

	/* sim_state: faulty/suspended -> critical */
	simSignal := diagnosisSignal{Key: signalKeySimState, State: signalStateOk,
		Detail: fmt.Sprintf("sim status is %q", in.SimStatus)}
	switch in.SimStatus {
	case "faulty", "suspended":
		simSignal.State = signalStateCritical
		simSignal.Detail = fmt.Sprintf("sim is %s", in.SimStatus)
	case "":
		simSignal.State = signalStateUnknown
		simSignal.Detail = "sim status unknown"
	}
	signals = append(signals, simSignal)

	/* package_state: expired -> warning */
	pkgSignal := diagnosisSignal{Key: signalKeyPackageState, State: signalStateOk,
		Detail: fmt.Sprintf("package status is %q", in.PackageStatus)}
	switch in.PackageStatus {
	case "expired":
		pkgSignal.State = signalStateWarning
		pkgSignal.Detail = "package has expired"
	case "":
		pkgSignal.State = signalStateUnknown
		pkgSignal.Detail = "package status unknown"
	}
	signals = append(signals, pkgSignal)

	/* site_health: uptime < 90 -> critical, < 99 -> warning */
	siteSignal := diagnosisSignal{Key: signalKeySiteHealth, State: signalStateUnknown,
		Detail: "no site health data"}
	if in.HasSiteHealth {
		switch {
		case in.SiteUptimePercent < 90:
			siteSignal.State = signalStateCritical
			siteSignal.Detail = fmt.Sprintf("site uptime is %.1f%%", in.SiteUptimePercent)
		case in.SiteUptimePercent < 99:
			siteSignal.State = signalStateWarning
			siteSignal.Detail = fmt.Sprintf("site uptime is %.1f%%", in.SiteUptimePercent)
		default:
			siteSignal.State = signalStateOk
			siteSignal.Detail = fmt.Sprintf("site uptime is %.1f%%", in.SiteUptimePercent)
		}
	}
	signals = append(signals, siteSignal)

	/* usage: zero usage over the last 24h -> warning */
	usageSignal := diagnosisSignal{Key: signalKeyUsage, State: signalStateOk,
		Detail: fmt.Sprintf("%.1f MB used in last 24h", in.UsageLast24hMb)}
	if in.UsageLast24hMb <= 0 {
		usageSignal.State = signalStateWarning
		usageSignal.Detail = "no data usage in last 24h"
	}
	signals = append(signals, usageSignal)

	/* last_seen: > 24h ago -> warning */
	lastSeenSignal := diagnosisSignal{Key: signalKeyLastSeen, State: signalStateUnknown,
		Detail: "never seen"}
	if in.LastSeenAt != nil {
		age := in.Now.Sub(*in.LastSeenAt)
		if age > 24*time.Hour {
			lastSeenSignal.State = signalStateWarning
			lastSeenSignal.Detail = fmt.Sprintf("last seen %.0f hours ago", age.Hours())
		} else {
			lastSeenSignal.State = signalStateOk
			lastSeenSignal.Detail = fmt.Sprintf("last seen %.0f minutes ago", age.Minutes())
		}
	}
	signals = append(signals, lastSeenSignal)

	/* likely issue = highest-severity signal; ordering above doubles as a
	   priority within the same severity (sim > package > site > usage > last_seen) */
	likely := pickLikelyIssue(signals)

	issue, action := describeIssue(likely, in)

	escalation := false
	for _, s := range signals {
		if s.State == signalStateCritical {
			escalation = true
			break
		}
	}

	return diagnosisResult{
		Signals:           signals,
		LikelyIssue:       issue,
		RecommendedAction: action,
		EscalationNeeded:  escalation,
	}
}

func severityRank(state string) int {
	switch state {
	case signalStateCritical:
		return 3
	case signalStateWarning:
		return 2
	case signalStateUnknown:
		return 1
	default: /* ok */
		return 0
	}
}

// pickLikelyIssue returns the first signal with the highest severity rank,
// or nil when everything is ok/unknown only.
func pickLikelyIssue(signals []diagnosisSignal) *diagnosisSignal {
	var best *diagnosisSignal
	bestRank := 0

	for i := range signals {
		r := severityRank(signals[i].State)
		if r >= 2 && r > bestRank {
			best = &signals[i]
			bestRank = r
		}
	}

	return best
}

// describeIssue maps the dominant signal to a human-readable likely issue
// and a recommended action.
func describeIssue(s *diagnosisSignal, in diagnosisInput) (issue, action string) {
	if s == nil {
		return "No issue detected", "No action needed"
	}

	switch s.Key {
	case signalKeySimState:
		if in.SimStatus == "faulty" {
			return "Faulty SIM", "Replace SIM"
		}

		return "SIM is suspended", "Review suspension and reactivate SIM"
	case signalKeyPackageState:
		return "Package expired", "Renew package"
	case signalKeySiteHealth:
		if s.State == signalStateCritical {
			return "Site is down or unstable", "Escalate to network team"
		}

		return "Site is degraded", "Escalate to network team"
	case signalKeyUsage:
		return "No recent data usage", "Check device"
	case signalKeyLastSeen:
		return "Customer not seen recently", "Check device"
	default:
		return "Unknown issue", "Check device"
	}
}
