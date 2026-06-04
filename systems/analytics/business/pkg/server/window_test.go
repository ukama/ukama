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
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ukama/ukama/systems/analytics/business/pb/gen"
)

// fixed reference: Wednesday 2026-06-03 15:30:00 UTC
var testNow = time.Date(2026, time.June, 3, 15, 30, 0, 0, time.UTC)

func TestResolveWindow_DefaultIsToday(t *testing.T) {
	win, err := resolveWindow(nil, testNow)
	assert.NoError(t, err)

	assert.Equal(t, PeriodToday, win.Period)
	assert.Equal(t, time.Date(2026, time.June, 3, 0, 0, 0, 0, time.UTC), win.From)
	assert.Equal(t, testNow, win.To)

	// previous window has the same length, ending where current starts
	assert.Equal(t, win.From, win.PrevTo)
	assert.Equal(t, win.To.Sub(win.From), win.PrevTo.Sub(win.PrevFrom))
}

func TestResolveWindow_Today(t *testing.T) {
	win, err := resolveWindow(&pb.Window{Period: PeriodToday}, testNow)
	assert.NoError(t, err)

	assert.Equal(t, time.Date(2026, time.June, 3, 0, 0, 0, 0, time.UTC), win.From)
	assert.Equal(t, testNow, win.To)
	assert.Equal(t, time.Date(2026, time.June, 2, 8, 30, 0, 0, time.UTC), win.PrevFrom)
	assert.Equal(t, time.Date(2026, time.June, 3, 0, 0, 0, 0, time.UTC), win.PrevTo)
}

func TestResolveWindow_Week(t *testing.T) {
	// 2026-06-03 is a Wednesday; ISO week starts Monday 2026-06-01.
	win, err := resolveWindow(&pb.Window{Period: PeriodWeek}, testNow)
	assert.NoError(t, err)

	assert.Equal(t, time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC), win.From)
	assert.Equal(t, testNow, win.To)
	assert.Equal(t, win.From, win.PrevTo)
}

func TestResolveWindow_WeekOnSunday(t *testing.T) {
	// 2026-06-07 is a Sunday; week still starts Monday 2026-06-01.
	sunday := time.Date(2026, time.June, 7, 10, 0, 0, 0, time.UTC)

	win, err := resolveWindow(&pb.Window{Period: PeriodWeek}, sunday)
	assert.NoError(t, err)

	assert.Equal(t, time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC), win.From)
}

func TestResolveWindow_Month(t *testing.T) {
	win, err := resolveWindow(&pb.Window{Period: PeriodMonth}, testNow)
	assert.NoError(t, err)

	assert.Equal(t, time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC), win.From)
	assert.Equal(t, testNow, win.To)
}

func TestResolveWindow_Custom(t *testing.T) {
	from := time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.May, 11, 0, 0, 0, 0, time.UTC)

	win, err := resolveWindow(&pb.Window{
		Period: PeriodCustom,
		From:   timestamppb.New(from),
		To:     timestamppb.New(to),
	}, testNow)
	assert.NoError(t, err)

	assert.True(t, win.From.Equal(from))
	assert.True(t, win.To.Equal(to))

	// prev window: 10 days immediately before from
	assert.True(t, win.PrevTo.Equal(from))
	assert.True(t, win.PrevFrom.Equal(from.AddDate(0, 0, -10)))
}

func TestResolveWindow_CustomMissingBounds(t *testing.T) {
	_, err := resolveWindow(&pb.Window{Period: PeriodCustom}, testNow)
	assert.Error(t, err)
}

func TestResolveWindow_CustomInvertedBounds(t *testing.T) {
	from := time.Date(2026, time.May, 11, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC)

	_, err := resolveWindow(&pb.Window{
		Period: PeriodCustom,
		From:   timestamppb.New(from),
		To:     timestamppb.New(to),
	}, testNow)
	assert.Error(t, err)
}

func TestResolveWindow_UnknownPeriod(t *testing.T) {
	_, err := resolveWindow(&pb.Window{Period: "fortnight"}, testNow)
	assert.Error(t, err)
}

func TestResolveWindow_Timezone(t *testing.T) {
	// 15:30 UTC on Jun 3 is 08:30 in America/Los_Angeles (UTC-7, PDT),
	// so "today" starts at Jun 3 00:00 PDT = Jun 3 07:00 UTC.
	win, err := resolveWindow(&pb.Window{
		Period:   PeriodToday,
		Timezone: "America/Los_Angeles",
	}, testNow)
	assert.NoError(t, err)

	expectedFrom := time.Date(2026, time.June, 3, 7, 0, 0, 0, time.UTC)
	assert.True(t, win.From.Equal(expectedFrom),
		"expected %v, got %v", expectedFrom, win.From)
}

func TestResolveWindow_InvalidTimezone(t *testing.T) {
	_, err := resolveWindow(&pb.Window{
		Period:   PeriodToday,
		Timezone: "Not/AZone",
	}, testNow)
	assert.Error(t, err)
}
