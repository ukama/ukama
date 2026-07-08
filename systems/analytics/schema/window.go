/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package schema

import (
	"time"
)

// Window is the atomic unit of the analytics pipeline. Windows are
// epoch-aligned: window N covers [N*W, (N+1)*W) in UTC. Window IDs are
// deterministic on any host, after any restart or replay — they are never
// derived from timers or process start.
type Window struct {
	ID    int64
	Start time.Time
	End   time.Time
}

// Grid is the shared window-grid configuration. It must be identical for
// ingest, analysis and aggregator (env-injected from the same deployment
// config).
type Grid struct {
	// W is the base window duration (pipeline default: 5m).
	W time.Duration
	// L is the watermark lag: a windowed pull for window N only becomes
	// eligible at N.End + L (pipeline default: 10m). Snapshot pulls are
	// eligible at N.End (state is captured at window close).
	L time.Duration
}

// WindowAt returns the window containing t.
func (g Grid) WindowAt(t time.Time) Window {
	sec := int64(g.W.Seconds())
	id := t.UTC().Unix() / sec

	return g.Window(id)
}

// Window returns the window with the given id.
func (g Grid) Window(id int64) Window {
	sec := int64(g.W.Seconds())

	return Window{
		ID:    id,
		Start: time.Unix(id*sec, 0).UTC(),
		End:   time.Unix((id+1)*sec, 0).UTC(),
	}
}

// EligibleAt returns the time at which the window may be pulled for the
// given strategy.
func (g Grid) EligibleAt(strategy Strategy, w Window) time.Time {
	if strategy == StrategyWindow {
		return w.End.Add(g.L)
	}

	return w.End
}

// NewestEligible returns the id of the most recent window eligible at 'now'
// for the given strategy: window N is eligible when now >= EligibleAt(N).
// Windowed pulls wait out the watermark L; snapshots are eligible at close.
func (g Grid) NewestEligible(strategy Strategy, now time.Time) int64 {
	ref := now
	if strategy == StrategyWindow {
		ref = now.Add(-g.L)
	}

	// The window containing ref has not closed yet; its predecessor is the
	// newest one whose eligibility time has passed.
	return g.WindowAt(ref).ID - 1
}
