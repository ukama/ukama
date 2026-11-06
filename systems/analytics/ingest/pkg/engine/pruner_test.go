/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package engine_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/systems/analytics/ingest/mocks"
	"github.com/ukama/ukama/systems/analytics/ingest/pkg/engine"
	"github.com/ukama/ukama/systems/analytics/schema"
)

const testOrg = "ukama"

var pruneSpecs = []schema.SourceSpec{{
	Source: "metrics",
	Pulls: []schema.PullSpec{
		{Key: "metrics.com_uptime.last", Strategy: schema.StrategyWindow},
		{Key: "metrics.active_sims.last", Strategy: schema.StrategyFullSnapshot},
		{Key: "metrics.ctl_uptime.last", Strategy: schema.StrategyWindow},
	},
}}

func TestPrunerManagesOnlyWindowDatasets(t *testing.T) {
	p := engine.NewPruner(schema.Grid{W: time.Minute}, pruneSpecs, &mocks.RawRepo{}, &mocks.LedgerRepo{},
		testOrg, 2880, time.Minute, 100)

	assert.Equal(t, []string{"metrics.com_uptime.last", "metrics.ctl_uptime.last"}, p.Datasets())
}

func TestPruneOnceDrainsBatchesAndLedger(t *testing.T) {
	grid := schema.Grid{W: time.Minute}
	now := time.Date(2026, 10, 29, 8, 0, 30, 0, time.UTC)
	const keep, batch = 2880, 2
	before := grid.NewestEligible(now) - keep + 1

	raw := &mocks.RawRepo{}
	ledger := &mocks.LedgerRepo{}

	raw.On("PruneBefore", testOrg, "metrics.com_uptime.last", before, batch).Return(int64(2), nil).Twice()
	raw.On("PruneBefore", testOrg, "metrics.com_uptime.last", before, batch).Return(int64(1), nil).Once()
	ledger.On("PruneBefore", testOrg, schema.LedgerKindDataset, "metrics.com_uptime.last", before, batch).
		Return(int64(0), nil).Once()

	raw.On("PruneBefore", testOrg, "metrics.ctl_uptime.last", before, batch).Return(int64(0), nil).Once()
	ledger.On("PruneBefore", testOrg, schema.LedgerKindDataset, "metrics.ctl_uptime.last", before, batch).
		Return(int64(2), nil).Once()
	ledger.On("PruneBefore", testOrg, schema.LedgerKindDataset, "metrics.ctl_uptime.last", before, batch).
		Return(int64(0), nil).Once()

	engine.NewPruner(grid, pruneSpecs, raw, ledger, testOrg, keep, time.Minute, batch).PruneOnce(now)

	raw.AssertExpectations(t)
	ledger.AssertExpectations(t)
}

func TestPruneOnceKeepsLedgerWhenRowDeleteFails(t *testing.T) {
	grid := schema.Grid{W: time.Minute}
	now := time.Date(2026, 10, 29, 8, 0, 30, 0, time.UTC)
	before := grid.NewestEligible(now) - 10 + 1

	raw := &mocks.RawRepo{}
	ledger := &mocks.LedgerRepo{}

	raw.On("PruneBefore", testOrg, "metrics.com_uptime.last", before, 100).
		Return(int64(0), errors.New("too many clients")).Once()
	raw.On("PruneBefore", testOrg, "metrics.ctl_uptime.last", before, 100).Return(int64(0), nil).Once()
	ledger.On("PruneBefore", testOrg, schema.LedgerKindDataset, "metrics.ctl_uptime.last", before, 100).
		Return(int64(0), nil).Once()

	engine.NewPruner(grid, pruneSpecs, raw, ledger, testOrg, 10, time.Minute, 100).PruneOnce(now)

	raw.AssertExpectations(t)
	ledger.AssertExpectations(t)
	ledger.AssertNotCalled(t, "PruneBefore", testOrg, schema.LedgerKindDataset, "metrics.com_uptime.last", before, 100)
}
