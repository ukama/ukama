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

	pb "github.com/ukama/ukama/systems/analytics/business/pb/gen"
)

const (
	PeriodToday  = "today"
	PeriodWeek   = "week"
	PeriodMonth  = "month"
	PeriodCustom = "custom"
)

// ResolvedWindow is a concrete [From, To) interval plus the previous window
// of equal length [PrevFrom, PrevTo) used for delta computation.
type ResolvedWindow struct {
	Period   string
	From     time.Time
	To       time.Time
	PrevFrom time.Time
	PrevTo   time.Time
	Location *time.Location
}

// resolveWindow turns a pb.Window into concrete boundaries. now is injected
// for testability. Defaults: nil window => "today" in UTC.
//
// Semantics (all half-open [from, to)):
//   - today: start of current day -> now
//   - week:  start of current ISO week (Monday) -> now
//   - month: start of current calendar month -> now
//   - custom: uses window.from/window.to verbatim (both required)
//
// The previous window is the interval of identical length immediately
// preceding From.
func resolveWindow(w *pb.Window, now time.Time) (*ResolvedWindow, error) {
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
		to = now
	case PeriodWeek:
		// ISO week starts on Monday.
		weekday := int(now.Weekday())
		if weekday == 0 { // Sunday
			weekday = 7
		}
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		from = startOfDay.AddDate(0, 0, -(weekday - 1))
		to = now
	case PeriodMonth:
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		to = now
	case PeriodCustom:
		if w == nil || w.From == nil || w.To == nil {
			return nil, fmt.Errorf("custom window requires both from and to")
		}
		from = w.From.AsTime().In(loc)
		to = w.To.AsTime().In(loc)
		if !to.After(from) {
			return nil, fmt.Errorf("custom window: to must be after from")
		}
	default:
		return nil, fmt.Errorf("unknown window period %q", period)
	}

	length := to.Sub(from)

	return &ResolvedWindow{
		Period:   period,
		From:     from,
		To:       to,
		PrevFrom: from.Add(-length),
		PrevTo:   from,
		Location: loc,
	}, nil
}
