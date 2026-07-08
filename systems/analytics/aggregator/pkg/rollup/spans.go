/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rollup

import (
	"fmt"
	"time"
)

const (
	SpanDaily   = "daily"
	SpanWeekly  = "weekly"
	SpanMonthly = "monthly"
)

// Spans is the configured span list (extensible to quarterly/annual without
// engine changes — add a case here and to the config list).
var Spans = []string{SpanDaily, SpanWeekly, SpanMonthly}

// SpanStart truncates t to the start of its span in loc: daily = midnight,
// weekly = ISO week (Monday), monthly = first of month.
func SpanStart(span string, t time.Time, loc *time.Location) (time.Time, error) {
	lt := t.In(loc)
	day := time.Date(lt.Year(), lt.Month(), lt.Day(), 0, 0, 0, 0, loc)

	switch span {
	case SpanDaily:
		return day, nil
	case SpanWeekly:
		offset := (int(day.Weekday()) + 6) % 7 // Monday = 0

		return day.AddDate(0, 0, -offset), nil
	case SpanMonthly:
		return time.Date(lt.Year(), lt.Month(), 1, 0, 0, 0, 0, loc), nil
	default:
		return time.Time{}, fmt.Errorf("unknown span %q", span)
	}
}

// SpanEnd returns the exclusive end of the span starting at start.
func SpanEnd(span string, start time.Time) (time.Time, error) {
	switch span {
	case SpanDaily:
		return start.AddDate(0, 0, 1), nil
	case SpanWeekly:
		return start.AddDate(0, 0, 7), nil
	case SpanMonthly:
		return start.AddDate(0, 1, 0), nil
	default:
		return time.Time{}, fmt.Errorf("unknown span %q", span)
	}
}

// PrevSpanStart returns the start of the span immediately before start.
func PrevSpanStart(span string, start time.Time) (time.Time, error) {
	switch span {
	case SpanDaily:
		return start.AddDate(0, 0, -1), nil
	case SpanWeekly:
		return start.AddDate(0, 0, -7), nil
	case SpanMonthly:
		return start.AddDate(0, -1, 0), nil
	default:
		return time.Time{}, fmt.Errorf("unknown span %q", span)
	}
}
