/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
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
	// W is the base window duration (pipeline default: 5m). Every window is
	// pulled as soon as it closes; late-arriving source data is handled by
	// the dirty-window recompute path, not by delaying pulls.
	W time.Duration
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

// NewestEligible returns the id of the most recent window pullable at 'now':
// a window is eligible as soon as it closes, so that is the predecessor of
// the window containing now.
func (g Grid) NewestEligible(now time.Time) int64 {
	return g.WindowAt(now).ID - 1
}

// MinCatchupWindows floors the catch-up horizon.
const MinCatchupWindows = 12

// CatchupWindows converts a catch-up duration into a window count, so the
// horizon is a fixed amount of time regardless of W. override > 0 pins an
// explicit count instead.
func CatchupWindows(catchup, w time.Duration, override int64) int64 {
	if override > 0 {
		return override
	}

	if w <= 0 || catchup <= 0 {
		return MinCatchupWindows
	}

	n := int64(catchup / w)
	if n < MinCatchupWindows {
		n = MinCatchupWindows
	}

	return n
}
