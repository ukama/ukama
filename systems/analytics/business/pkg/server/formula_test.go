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

	"github.com/ukama/ukama/systems/analytics/business/pkg/db"
	"github.com/ukama/ukama/systems/common/uuid"
)

func TestPctDelta(t *testing.T) {
	assert.Equal(t, 0.0, pctDelta(100, 0), "zero baseline guards to 0")
	assert.Equal(t, 0.0, pctDelta(0, 0))
	assert.InDelta(t, 100.0, pctDelta(200, 100), 1e-9)
	assert.InDelta(t, -50.0, pctDelta(50, 100), 1e-9)
	assert.InDelta(t, -100.0, pctDelta(0, 100), 1e-9)
}

func TestAvgPurchase(t *testing.T) {
	assert.Equal(t, 0.0, avgPurchase(100, 0), "zero purchases guards to 0")
	assert.InDelta(t, 25.0, avgPurchase(100, 4), 1e-9)
}

func TestMrr(t *testing.T) {
	assert.Equal(t, 0.0, mrr(nil))

	pkgs := []db.PackageSnapshot{
		{Price: 10, ActiveSubscribers: 5},  // 50
		{Price: 20, ActiveSubscribers: 0},  // 0
		{Price: 2.5, ActiveSubscribers: 4}, // 10
	}

	assert.InDelta(t, 60.0, mrr(pkgs), 1e-9)
}

func TestActiveSubscribers(t *testing.T) {
	assert.Equal(t, uint32(0), activeSubscribers(nil))

	pkgs := []db.PackageSnapshot{
		{ActiveSubscribers: 5},
		{ActiveSubscribers: 7},
	}

	assert.Equal(t, uint32(12), activeSubscribers(pkgs))
}

func TestArpu(t *testing.T) {
	assert.Equal(t, 0.0, arpu(100, 0), "zero active customers guards to 0")
	assert.InDelta(t, 10.0, arpu(100, 10), 1e-9)
}

func TestTopPlanByRevenue(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		_, _, _, ok := topPlanByRevenue(nil)
		assert.False(t, ok)
	})

	t.Run("AllZeroRevenue", func(t *testing.T) {
		rollups := []db.BusinessPackageRollupDaily{
			{PackageId: uuid.NewV4(), Revenue: 0},
		}

		_, _, _, ok := topPlanByRevenue(rollups)
		assert.False(t, ok)
	})

	t.Run("PicksTopAndComputesShare", func(t *testing.T) {
		day := time.Now()
		a := uuid.NewV4()
		b := uuid.NewV4()

		rollups := []db.BusinessPackageRollupDaily{
			{Day: day, PackageId: a, Revenue: 30},
			{Day: day.AddDate(0, 0, 1), PackageId: a, Revenue: 30},
			{Day: day, PackageId: b, Revenue: 40},
		}

		topId, topRevenue, share, ok := topPlanByRevenue(rollups)
		assert.True(t, ok)
		assert.Equal(t, a.String(), topId)
		assert.InDelta(t, 60.0, topRevenue, 1e-9)
		assert.InDelta(t, 60.0, share, 1e-9)
	})
}

func TestMakeKpiFormatting(t *testing.T) {
	now := time.Now()

	k := makeKpi("revenue", 1234.5, formatMoney(1234.5), 12.34, "month", now)
	assert.Equal(t, "revenue", k.Key)
	assert.Equal(t, 1234.5, k.Value)
	assert.Equal(t, "$1234.50", k.Formatted)
	assert.Equal(t, 12.34, k.Delta)
	assert.Equal(t, "month", k.DeltaPeriod)
	assert.False(t, k.Stale)
	assert.Equal(t, now.Unix(), k.AsOf.AsTime().Unix())

	assert.Equal(t, "42", formatCount(42))
	assert.Equal(t, "99.5%", formatPercent(99.5))
}

func TestMakeMeta(t *testing.T) {
	m := makeMeta(45, 0, 20)
	assert.Equal(t, uint32(45), m.Count)
	assert.Equal(t, uint32(1), m.Page, "page defaults to 1")
	assert.Equal(t, uint32(20), m.Size)
	assert.Equal(t, uint32(3), m.Pages)

	m = makeMeta(0, 2, 20)
	assert.Equal(t, uint32(0), m.Pages)

	m = makeMeta(10, 1, 0)
	assert.Equal(t, uint32(0), m.Pages, "zero size guards pages to 0")
}
