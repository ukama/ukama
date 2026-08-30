/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package schema

import (
	"testing"
	"time"
)

// The horizon stays a fixed amount of time as W changes.
func TestCatchupWindows_ScalesWithWindowSize(t *testing.T) {
	cases := []struct {
		name string
		w    time.Duration
		want int64
	}{
		{"five minute windows", 5 * time.Minute, 12},
		{"one minute windows", time.Minute, 60},
		{"ten second windows", 10 * time.Second, 360},
		{"one second windows", time.Second, 3600},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := CatchupWindows(time.Hour, c.w, 0); got != c.want {
				t.Errorf("CatchupWindows(1h, %s) = %d, want %d (= 1h of recovery)", c.w, got, c.want)
			}
		})
	}
}

func TestCatchupWindows_FlooredAndOverridable(t *testing.T) {
	// A window longer than the horizon still gets a usable minimum.
	if got := CatchupWindows(time.Hour, 24*time.Hour, 0); got != MinCatchupWindows {
		t.Errorf("oversized window: got %d, want the %d floor", got, MinCatchupWindows)
	}

	// Explicit override wins.
	if got := CatchupWindows(time.Hour, time.Minute, 5); got != 5 {
		t.Errorf("override: got %d, want 5", got)
	}

	// Degenerate config never returns 0 (which would scan nothing at all).
	if got := CatchupWindows(0, 0, 0); got != MinCatchupWindows {
		t.Errorf("zero config: got %d, want the %d floor", got, MinCatchupWindows)
	}
}
