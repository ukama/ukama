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

	pb "github.com/ukama/ukama/systems/analytics/customer/pb/gen"
)

const (
	PeriodToday  = "today"
	PeriodWeek   = "week"
	PeriodMonth  = "month"
	PeriodCustom = "custom"
)

// timeWindow is a resolved half-open interval [From, To) along with the
// previous window of equal length [PrevFrom, PrevTo) used for deltas.
type timeWindow struct {
	Period   string
	From     time.Time
	To       time.Time
	PrevFrom time.Time
	PrevTo   time.Time
	Location *time.Location
}

// resolveWindow converts the request window into concrete time bounds.
// A nil window defaults to "today" in UTC. For "custom", both from and to
// must be set. The previous window always has equal length and ends where
// the current window starts.
func resolveWindow(w *pb.Window, now time.Time) (*timeWindow, error) {
	loc := time.UTC

	period := PeriodToday
	if w != nil && w.Period != "" {
		period = w.Period
	}

	if w != nil && w.Timezone != "" {
		l, err := time.LoadLocation(w.Timezone)
		if err != nil {
			return nil, fmt.Errorf("invalid timezone %q: %w", w.Timezone, err)
		}
		loc = l
	}

	now = now.In(loc)

	var from, to time.Time

	switch period {
	case PeriodToday:
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		to = from.AddDate(0, 0, 1)
	case PeriodWeek:
		// week starts on Monday
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).
			AddDate(0, 0, -(weekday - 1))
		to = from.AddDate(0, 0, 7)
	case PeriodMonth:
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		to = from.AddDate(0, 1, 0)
	case PeriodCustom:
		if w == nil || w.From == nil || w.To == nil {
			return nil, fmt.Errorf("custom period requires both from and to")
		}

		from = w.From.AsTime().In(loc)
		to = w.To.AsTime().In(loc)

		if !to.After(from) {
			return nil, fmt.Errorf("window to must be after from")
		}
	default:
		return nil, fmt.Errorf("invalid period %q: must be one of today, week, month, custom", period)
	}

	length := to.Sub(from)

	return &timeWindow{
		Period:   period,
		From:     from,
		To:       to,
		PrevFrom: from.Add(-length),
		PrevTo:   from,
		Location: loc,
	}, nil
}
