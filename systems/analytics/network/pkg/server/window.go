/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"time"

	pb "github.com/ukama/ukama/systems/analytics/network/pb/gen"
)

// TimeWindow is a resolved time window with its previous (comparison) window.
type TimeWindow struct {
	Period   string
	From     time.Time
	To       time.Time
	PrevFrom time.Time
	PrevTo   time.Time
}

// ResolveWindow turns a request Window into concrete [from, to) bounds plus
// the immediately preceding window of the same length (for deltas).
// Supported periods: today, week, month, custom (uses from/to from the
// request). Defaults to today when the window is nil or unknown.
func ResolveWindow(w *pb.Window, now time.Time) TimeWindow {
	loc := time.UTC

	if w != nil && w.Timezone != "" {
		if l, err := time.LoadLocation(w.Timezone); err == nil {
			loc = l
		}
	}

	now = now.In(loc)

	period := "today"
	if w != nil && w.Period != "" {
		period = w.Period
	}

	var from, to time.Time

	switch {
	case period == "week":
		to = now
		from = now.AddDate(0, 0, -7)
	case period == "month":
		to = now
		from = now.AddDate(0, -1, 0)
	case period == "custom" && w != nil && w.From != nil && w.To != nil:
		from = w.From.AsTime().In(loc)
		to = w.To.AsTime().In(loc)
	default:
		period = "today"
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		to = now
	}

	length := to.Sub(from)

	return TimeWindow{
		Period:   period,
		From:     from,
		To:       to,
		PrevFrom: from.Add(-length),
		PrevTo:   from,
	}
}
