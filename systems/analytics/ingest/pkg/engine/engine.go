/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package engine

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/analytics/ingest/pkg/db"
	"github.com/ukama/ukama/systems/analytics/schema"
	"github.com/ukama/ukama/systems/common/msgbus"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
)

// Engine is the ingest scheduler: every tick it walks the pull DAG in
// topological stages and, for each dataset, claims and executes all eligible,
// not-yet-pulled windows. The ledger claim is what prevents double-pulling —
// across ticks, restarts, replays and replicas.
type Engine struct {
	grid    schema.Grid
	stages  [][]schema.PullSpec
	ledger  db.LedgerRepo
	errors  db.ErrorRepo
	puller  *Puller
	mb      mb.MsgBusServiceClient
	org     string
	tick    time.Duration
	catchup int64

	windowReadyRoute string
	stop             chan struct{}
}

func New(grid schema.Grid, specs []schema.SourceSpec, ledger db.LedgerRepo,
	errors db.ErrorRepo, puller *Puller, msgBus mb.MsgBusServiceClient,
	org string, tick time.Duration, catchup int64) (*Engine, error) {
	stages, err := schema.OrderPulls(specs)
	if err != nil {
		return nil, err
	}

	return &Engine{
		grid:             grid,
		stages:           stages,
		ledger:           ledger,
		errors:           errors,
		puller:           puller,
		mb:               msgBus,
		org:              org,
		tick:             tick,
		catchup:          catchup,
		windowReadyRoute: msgbus.PrepareRoute(org, schema.WindowReadyRoute),
		stop:             make(chan struct{}),
	}, nil
}

func (e *Engine) Start() {
	log.Infof("ingest engine starting: W-grid ticks every %s, %d pull stages", e.tick, len(e.stages))

	ticker := time.NewTicker(e.tick)
	defer ticker.Stop()

	e.runOnce() // immediate first pass

	for {
		select {
		case <-ticker.C:
			e.runOnce()
		case <-e.stop:
			return
		}
	}
}

func (e *Engine) Stop() {
	close(e.stop)
}

// runOnce processes stages in DAG order so for_each children see their
// parents' data for the same window.
func (e *Engine) runOnce() {
	now := time.Now().UTC()

	for _, stage := range e.stages {
		for _, pull := range stage {
			if err := e.processPull(pull, now); err != nil {
				log.Errorf("processing dataset %s: %v", pull.Key, err)
			}
		}
	}
}

// processPull scans the whole catchup range every tick (not just past the
// newest pulled window) so failed or crashed-in-flight windows behind the
// high-water mark are retried too. Cheap: completed windows short-circuit on
// a ledger status read.
func (e *Engine) processPull(pull schema.PullSpec, now time.Time) error {
	newest := e.grid.NewestEligible(now)

	for w := newest - e.catchup + 1; w <= newest; w++ {
		status, err := e.ledger.Status(e.org, schema.LedgerKindDataset, pull.Key, w)
		if err != nil {
			return err
		}

		if status == schema.StatusPulled {
			continue
		}

		if err := e.processWindow(pull, w); err != nil {
			log.Errorf("dataset %s window %d: %v", pull.Key, w, err)
		}
	}

	return nil
}

func (e *Engine) processWindow(pull schema.PullSpec, windowID int64) error {
	// A for_each child needs its parent pulled for the same window first.
	if pull.ForEach != nil {
		parentStatus, err := e.ledger.Status(e.org, schema.LedgerKindDataset, pull.ForEach.Dataset, windowID)
		if err != nil {
			return err
		}

		if parentStatus != schema.StatusPulled {
			return nil // parent pending; retried next tick
		}
	}

	claimed, err := e.ledger.Claim(e.org, schema.LedgerKindDataset, pull.Key, windowID)
	if err != nil {
		return err
	}

	if !claimed {
		return nil // already pulled, in flight elsewhere, or terminally done
	}

	win := e.grid.Window(windowID)

	n, err := e.puller.Execute(pull, win)
	if err != nil {
		if lerr := e.ledger.Mark(e.org, schema.LedgerKindDataset, pull.Key, windowID,
			schema.StatusFailed, err.Error()); lerr != nil {
			log.Errorf("marking ledger failed: %v", lerr)
		}

		if derr := e.errors.Insert(e.org, pull.Key, windowID, "", err.Error()); derr != nil {
			log.Errorf("recording ingest error: %v", derr)
		}

		return err
	}

	if err := e.ledger.Mark(e.org, schema.LedgerKindDataset, pull.Key, windowID,
		schema.StatusPulled, fmt.Sprintf("records=%d", n)); err != nil {
		return err
	}

	log.Infof("dataset %s window %d pulled (%d records)", pull.Key, windowID, n)

	e.publishWindowReady(pull.Key, windowID)

	return nil
}

// publishWindowReady is the fast path only: analysis' sweeper recovers from
// lost events via the ledger.
func (e *Engine) publishWindowReady(datasetKey string, windowID int64) {
	payload, err := schema.WindowReady{
		DatasetKey: datasetKey,
		WindowID:   windowID,
		OrgName:    e.org,
	}.ToStruct()
	if err != nil {
		log.Errorf("building window.ready payload: %v", err)

		return
	}

	if err := e.mb.PublishRequest(e.windowReadyRoute, payload); err != nil {
		log.Errorf("publishing window.ready for %s window %d: %v", datasetKey, windowID, err)
	}
}
