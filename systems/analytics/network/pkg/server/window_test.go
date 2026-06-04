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

	"github.com/tj/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ukama/ukama/systems/analytics/network/pb/gen"
)

func TestResolveWindow(t *testing.T) {
	now := time.Date(2026, 6, 3, 15, 30, 0, 0, time.UTC)

	t.Run("nil window defaults to today", func(t *testing.T) {
		w := ResolveWindow(nil, now)

		assert.Equal(t, "today", w.Period)
		assert.Equal(t, time.Date(2026, 6, 3, 0, 0, 0, 0, time.UTC), w.From)
		assert.Equal(t, now, w.To)
		assert.Equal(t, w.From, w.PrevTo)
		assert.Equal(t, w.To.Sub(w.From), w.PrevTo.Sub(w.PrevFrom))
	})

	t.Run("week", func(t *testing.T) {
		w := ResolveWindow(&pb.Window{Period: "week"}, now)

		assert.Equal(t, "week", w.Period)
		assert.Equal(t, now.AddDate(0, 0, -7), w.From)
		assert.Equal(t, now, w.To)
		assert.Equal(t, w.From, w.PrevTo)
	})

	t.Run("month", func(t *testing.T) {
		w := ResolveWindow(&pb.Window{Period: "month"}, now)

		assert.Equal(t, "month", w.Period)
		assert.Equal(t, now.AddDate(0, -1, 0), w.From)
		assert.Equal(t, now, w.To)
	})

	t.Run("custom", func(t *testing.T) {
		from := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)

		w := ResolveWindow(&pb.Window{
			Period: "custom",
			From:   timestamppb.New(from),
			To:     timestamppb.New(to),
		}, now)

		assert.Equal(t, "custom", w.Period)
		assert.Equal(t, from, w.From)
		assert.Equal(t, to, w.To)
		assert.Equal(t, from.AddDate(0, 0, -1), w.PrevFrom)
		assert.Equal(t, from, w.PrevTo)
	})

	t.Run("custom without bounds falls back to today", func(t *testing.T) {
		w := ResolveWindow(&pb.Window{Period: "custom"}, now)

		assert.Equal(t, "today", w.Period)
	})

	t.Run("unknown period falls back to today", func(t *testing.T) {
		w := ResolveWindow(&pb.Window{Period: "fortnight"}, now)

		assert.Equal(t, "today", w.Period)
	})
}
