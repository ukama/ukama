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

	pb "github.com/ukama/ukama/systems/analytics/customer/pb/gen"
)

func TestResolveWindow(t *testing.T) {
	// Wednesday, June 3, 2026, 15:30 UTC
	now := time.Date(2026, 6, 3, 15, 30, 0, 0, time.UTC)

	t.Run("nil window defaults to today", func(t *testing.T) {
		w, err := resolveWindow(nil, now)
		assert.NoError(t, err)
		assert.Equal(t, PeriodToday, w.Period)
		assert.Equal(t, time.Date(2026, 6, 3, 0, 0, 0, 0, time.UTC), w.From)
		assert.Equal(t, time.Date(2026, 6, 4, 0, 0, 0, 0, time.UTC), w.To)
		assert.Equal(t, time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC), w.PrevFrom)
		assert.Equal(t, w.From, w.PrevTo)
	})

	t.Run("week starts on monday", func(t *testing.T) {
		w, err := resolveWindow(&pb.Window{Period: PeriodWeek}, now)
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), w.From)
		assert.Equal(t, time.Date(2026, 6, 8, 0, 0, 0, 0, time.UTC), w.To)
		assert.Equal(t, time.Date(2026, 5, 25, 0, 0, 0, 0, time.UTC), w.PrevFrom)
	})

	t.Run("week handles sunday as last day", func(t *testing.T) {
		sunday := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
		w, err := resolveWindow(&pb.Window{Period: PeriodWeek}, sunday)
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), w.From)
	})

	t.Run("month window", func(t *testing.T) {
		w, err := resolveWindow(&pb.Window{Period: PeriodMonth}, now)
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), w.From)
		assert.Equal(t, time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC), w.To)
		assert.Equal(t, time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC), w.PrevFrom)
	})

	t.Run("custom window with prev of equal length", func(t *testing.T) {
		from := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2026, 6, 3, 0, 0, 0, 0, time.UTC)

		w, err := resolveWindow(&pb.Window{
			Period: PeriodCustom,
			From:   timestamppb.New(from),
			To:     timestamppb.New(to),
		}, now)
		assert.NoError(t, err)
		assert.Equal(t, from, w.From)
		assert.Equal(t, to, w.To)
		assert.Equal(t, time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC), w.PrevFrom)
		assert.Equal(t, from, w.PrevTo)
	})

	t.Run("custom without bounds fails", func(t *testing.T) {
		_, err := resolveWindow(&pb.Window{Period: PeriodCustom}, now)
		assert.Error(t, err)
	})

	t.Run("custom with inverted bounds fails", func(t *testing.T) {
		from := time.Date(2026, 6, 3, 0, 0, 0, 0, time.UTC)
		to := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

		_, err := resolveWindow(&pb.Window{
			Period: PeriodCustom,
			From:   timestamppb.New(from),
			To:     timestamppb.New(to),
		}, now)
		assert.Error(t, err)
	})

	t.Run("invalid period fails", func(t *testing.T) {
		_, err := resolveWindow(&pb.Window{Period: "year"}, now)
		assert.Error(t, err)
	})

	t.Run("invalid timezone fails", func(t *testing.T) {
		_, err := resolveWindow(&pb.Window{Period: PeriodToday, Timezone: "Not/AZone"}, now)
		assert.Error(t, err)
	})

	t.Run("timezone shifts the day boundary", func(t *testing.T) {
		w, err := resolveWindow(&pb.Window{Period: PeriodToday, Timezone: "America/Los_Angeles"}, now)
		assert.NoError(t, err)

		loc, _ := time.LoadLocation("America/Los_Angeles")
		assert.Equal(t, time.Date(2026, 6, 3, 0, 0, 0, 0, loc), w.From)
	})
}
