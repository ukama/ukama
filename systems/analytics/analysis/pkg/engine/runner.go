/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package engine

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/analytics/analysis/pkg/algos"
	"github.com/ukama/ukama/systems/analytics/analysis/pkg/db"
	"github.com/ukama/ukama/systems/analytics/schema"
	"github.com/ukama/ukama/systems/common/msgbus"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
)

// Runner computes KPIs: triggered by window.ready events (fast path) and a
// ledger sweeper (recovery path). A KPI computes for a window once ALL its
// input datasets are pulled for that window.
type Runner struct {
	grid     schema.Grid
	kpis     []schema.KpiSpec
	byInput  map[string][]schema.KpiSpec // dataset key -> dependent KPI specs
	registry *algos.Registry

	raw    db.RawReader
	kpiDb  db.KpiRepo
	ledger db.LedgerRepo
	errs   db.ErrorRepo
	mb     mb.MsgBusServiceClient

	org              string
	sweep            time.Duration
	catchup          int64
	kpiComputedRoute string
	stop             chan struct{}
}

func NewRunner(grid schema.Grid, kpis []schema.KpiSpec, registry *algos.Registry,
	raw db.RawReader, kpiDb db.KpiRepo, ledger db.LedgerRepo, errs db.ErrorRepo,
	msgBus mb.MsgBusServiceClient, org string, sweep time.Duration, catchup int64) (*Runner, error) {
	// Fail fast: every spec's algo must exist in the registry.
	for _, k := range kpis {
		if _, err := registry.Get(k.Algo); err != nil {
			return nil, fmt.Errorf("kpi %s: %w", k.Kpi, err)
		}
	}

	byInput := map[string][]schema.KpiSpec{}
	for _, k := range kpis {
		for _, ds := range k.InputDatasets() {
			byInput[ds] = append(byInput[ds], k)
		}
	}

	return &Runner{
		grid:             grid,
		kpis:             kpis,
		byInput:          byInput,
		registry:         registry,
		raw:              raw,
		kpiDb:            kpiDb,
		ledger:           ledger,
		errs:             errs,
		mb:               msgBus,
		org:              org,
		sweep:            sweep,
		catchup:          catchup,
		kpiComputedRoute: msgbus.PrepareRoute(org, schema.KpiComputedRoute),
		stop:             make(chan struct{}),
	}, nil
}

// OnDatasetReady is the event fast path (called by the msgbus event server).
func (r *Runner) OnDatasetReady(datasetKey string, windowID int64) {
	for _, kpi := range r.byInput[datasetKey] {
		if err := r.TryCompute(kpi, windowID); err != nil {
			log.Errorf("kpi %s window %d: %v", kpi.Kpi, windowID, err)
		}
	}
}

// StartSweeper recovers windows whose events were lost; the ledger is the
// source of truth, events are just latency.
func (r *Runner) StartSweeper() {
	ticker := time.NewTicker(r.sweep)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.sweepOnce()
		case <-r.stop:
			return
		}
	}
}

func (r *Runner) Stop() {
	close(r.stop)
}

func (r *Runner) sweepOnce() {
	now := time.Now().UTC()
	newest := r.grid.NewestEligible(now)

	for _, kpi := range r.kpis {
		for w := newest - r.catchup + 1; w <= newest; w++ {
			status, err := r.ledger.KpiStatus(r.org, kpi.Kpi, w)
			if err != nil {
				log.Errorf("sweeper: kpi %s window %d status: %v", kpi.Kpi, w, err)

				continue
			}

			if status == schema.StatusComputed || status == schema.StatusInFlight {
				continue
			}

			if err := r.TryCompute(kpi, w); err != nil {
				log.Errorf("sweeper: kpi %s window %d: %v", kpi.Kpi, w, err)
			}
		}
	}
}

// TryCompute computes one KPI for one window if all inputs are ready.
// Idempotent: recomputation upserts identical rows.
func (r *Runner) TryCompute(kpi schema.KpiSpec, windowID int64) error {
	ready, err := r.inputsReady(kpi, windowID)
	if err != nil {
		return err
	}

	if !ready {
		return nil // retried on next event/sweep
	}

	claimed, err := r.ledger.ClaimKpi(r.org, kpi.Kpi, windowID)
	if err != nil {
		return err
	}

	if !claimed {
		return nil
	}

	if err := r.compute(kpi, windowID); err != nil {
		if lerr := r.ledger.MarkKpi(r.org, kpi.Kpi, windowID, schema.StatusFailed, err.Error()); lerr != nil {
			log.Errorf("marking kpi ledger failed: %v", lerr)
		}

		if derr := r.errs.Insert(r.org, kpi.Kpi, windowID, err.Error()); derr != nil {
			log.Errorf("recording analysis error: %v", derr)
		}

		return err
	}

	return r.ledger.MarkKpi(r.org, kpi.Kpi, windowID, schema.StatusComputed, "")
}

func (r *Runner) inputsReady(kpi schema.KpiSpec, windowID int64) (bool, error) {
	for _, ds := range kpi.InputDatasets() {
		status, err := r.ledger.DatasetStatus(r.org, ds, windowID)
		if err != nil {
			return false, err
		}

		if status != schema.StatusPulled {
			return false, nil
		}
	}

	return true, nil
}

func (r *Runner) compute(kpi schema.KpiSpec, windowID int64) error {
	win := r.grid.Window(windowID)

	inputs, err := r.loadInputs(kpi, windowID)
	if err != nil {
		return err
	}

	algo, err := r.registry.Get(kpi.Algo)
	if err != nil {
		return err
	}

	results, err := algo(win, inputs, kpi)
	if err != nil {
		return fmt.Errorf("algo %s: %w", kpi.Algo, err)
	}

	now := time.Now().UTC()
	rows := make([]schema.KpiWindow, 0, len(results))

	for _, res := range results {
		rows = append(rows, schema.KpiWindow{
			KpiKey:      kpi.Kpi,
			OrgID:       r.org,
			Scope:       schema.CanonicalScope(res.Scope),
			WindowID:    windowID,
			Value:       res.Value,
			Sum:         res.Sum,
			Count:       res.Count,
			Min:         res.Min,
			Max:         res.Max,
			ValueType:   kpi.Output.Type,
			Unit:        kpi.Output.Unit,
			Symbol:      kpi.Output.Symbol,
			AlgoVersion: kpi.Algo,
			ComputedAt:  now,
		})
	}

	// Whole window per KPI is one write; no partial silent writes.
	if err := r.kpiDb.Upsert(rows); err != nil {
		return err
	}

	log.Infof("kpi %s window %d computed (%d scope rows)", kpi.Kpi, windowID, len(rows))

	r.publishKpiComputed(kpi.Kpi, windowID)

	return nil
}

func (r *Runner) loadInputs(kpi schema.KpiSpec, windowID int64) (algos.Datasets, error) {
	inputs := algos.Datasets{}

	for name, in := range kpi.Inputs {
		var records []schema.RawRecord

		var err error

		if in.Mode == "window" {
			records, err = r.raw.WindowRows(r.org, in.Dataset, windowID)
		} else {
			records, err = r.raw.StateAsOf(r.org, in.Dataset, windowID)
		}

		if err != nil {
			return nil, fmt.Errorf("loading input %s (%s): %w", name, in.Dataset, err)
		}

		rows := make([]map[string]interface{}, 0, len(records))

		for _, rec := range records {
			fields := map[string]interface{}{}
			if err := json.Unmarshal([]byte(rec.Fields), &fields); err != nil {
				return nil, fmt.Errorf("input %s record %d fields: %w", name, rec.ID, err)
			}

			rows = append(rows, fields)
		}

		inputs[name] = rows
	}

	return inputs, nil
}

func (r *Runner) publishKpiComputed(kpiKey string, windowID int64) {
	payload, err := schema.KpiComputed{
		KpiKey:   kpiKey,
		WindowID: windowID,
		OrgName:  r.org,
	}.ToStruct()
	if err != nil {
		log.Errorf("building kpi.computed payload: %v", err)

		return
	}

	if err := r.mb.PublishRequest(r.kpiComputedRoute, payload); err != nil {
		log.Errorf("publishing kpi.computed for %s window %d: %v", kpiKey, windowID, err)
	}
}
