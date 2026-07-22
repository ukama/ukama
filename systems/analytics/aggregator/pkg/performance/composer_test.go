/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package performance

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/systems/analytics/schema"
)

func day(d int) time.Time { return time.Date(2026, 7, 1+d, 0, 0, 0, 0, time.UTC) }

func rollup(net, pkg string, value float64, spanStart time.Time) schema.KpiRollup {
	return schema.KpiRollup{
		Scope:     schema.CanonicalScope(map[string]string{"network_id": net, "package_id": pkg}),
		Value:     value,
		SpanStart: spanStart,
	}
}

func TestGroupByEntity(t *testing.T) {
	rows := []schema.KpiRollup{
		rollup("n1", "pkgA", 5, day(0)),
		rollup("n2", "pkgA", 3, day(1)),
		rollup("n1", "pkgB", 100, day(2)),
	}

	byEntity := groupByEntity(rows, "package_id")
	assert.Len(t, byEntity["pkgA"], 2, "pkgA rows across two networks")
	assert.Len(t, byEntity["pkgB"], 1)
}

// SUM folds every daily row of an entity across days AND scopes — this is the
// whole point of the rolling report window (stable totals vs a single day).
func TestFoldValue_SumAcrossWindow(t *testing.T) {
	rows := []schema.KpiRollup{
		rollup("n1", "pkgA", 5, day(0)),
		rollup("n1", "pkgA", 7, day(1)),
		rollup("n2", "pkgA", 3, day(1)),
	}

	v, ok := foldValue("SUM", rows)
	assert.True(t, ok)
	assert.Equal(t, float64(15), v)

	mx, _ := foldValue("MAX", rows)
	assert.Equal(t, float64(7), mx)

	avg, _ := foldValue("AVG", rows)
	assert.Equal(t, float64(5), avg) // 15 / 3

	_, ok = foldValue("SUM", nil)
	assert.False(t, ok)
}

func TestFoldValue_LastByMostRecentDay(t *testing.T) {
	rows := []schema.KpiRollup{
		rollup("n1", "pkgA", 10, day(0)),
		rollup("n1", "pkgA", 42, day(3)), // most recent
		rollup("n1", "pkgA", 20, day(1)),
	}

	v, _ := foldValue("LAST", rows)
	assert.Equal(t, float64(42), v)
}

func TestComposeCell_Trend(t *testing.T) {
	col := schema.ReportColumn{Name: "sold", Op: "SUM"}
	rows := []schema.KpiRollup{rollup("n1", "pkgA", 15, day(1))}

	// no previous window -> new
	assert.Equal(t, "new", composeCell(col, rows, 0, false).Trend)

	// up vs previous window
	up := composeCell(col, rows, 10, true)
	assert.Equal(t, "up", up.Trend)
	assert.Equal(t, float64(15), up.Value)
	assert.NotNil(t, up.PrevValue)
	assert.Equal(t, float64(10), *up.PrevValue)
	assert.Equal(t, float64(5), *up.ChangeAbs)

	assert.Equal(t, "down", composeCell(col, rows, 20, true).Trend)
	assert.Equal(t, "na", composeCell(col, rows, 0, true).Trend) // prev 0, cur>0

	// empty entity -> zero cell, no data
	empty := composeCell(col, nil, 0, false)
	assert.False(t, empty.hasData)
	assert.Equal(t, float64(0), empty.Value)
}

func TestReportWindowLabel(t *testing.T) {
	assert.Equal(t, "8w", reportWindowLabel(1344*time.Hour))
	assert.Equal(t, "3d", reportWindowLabel(72*time.Hour))
}
