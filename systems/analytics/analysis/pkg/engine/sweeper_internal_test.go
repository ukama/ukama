/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package engine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/ukama/ukama/systems/analytics/analysis/mocks"
	"github.com/ukama/ukama/systems/analytics/schema"
)

func sweepRunner(ledger *mocks.LedgerRepo) *Runner {
	return &Runner{
		grid: schema.Grid{W: time.Minute},
		kpis: []schema.KpiSpec{{
			Kpi:    "REVENUE",
			Inputs: map[string]schema.InputSpec{"payments": {Dataset: "payments.processor.list"}},
		}},
		ledger:  ledger,
		org:     "ukama",
		catchup: 1,
		stop:    make(chan struct{}),
	}
}

func TestSweepRetriesInFlightWindowViaClaim(t *testing.T) {
	ledger := &mocks.LedgerRepo{}
	ledger.On("KpiStatus", "ukama", "REVENUE", mock.Anything).Return(schema.StatusInFlight, nil)
	ledger.On("DatasetStatus", "ukama", "payments.processor.list", mock.Anything).Return(schema.StatusPulled, nil)
	ledger.On("ClaimKpi", "ukama", "REVENUE", mock.Anything).Return(false, nil).Once()

	sweepRunner(ledger).sweepOnce()

	ledger.AssertExpectations(t)
}

func TestSweepSkipsComputedWindow(t *testing.T) {
	ledger := &mocks.LedgerRepo{}
	ledger.On("KpiStatus", "ukama", "REVENUE", mock.Anything).Return(schema.StatusComputed, nil)

	sweepRunner(ledger).sweepOnce()

	ledger.AssertExpectations(t)
	ledger.AssertNotCalled(t, "ClaimKpi", mock.Anything, mock.Anything, mock.Anything)
}
