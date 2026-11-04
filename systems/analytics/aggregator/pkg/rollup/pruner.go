/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rollup

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/db"
	"github.com/ukama/ukama/systems/analytics/schema"
)

// Pruner deletes kpi_windows rows and their KPI ledger entries once they are
// older than keep. It waits for ready (the engine's first clean backfill) so
// no window is deleted before its spans are rolled up.
type Pruner struct {
	grid     schema.Grid
	repo     db.RetentionRepo
	org      string
	keep     time.Duration
	interval time.Duration
	batch    int
	ready    <-chan struct{}
	stop     chan struct{}
	done     chan struct{}
}

func NewPruner(grid schema.Grid, repo db.RetentionRepo, org string, keep, interval time.Duration,
	batch int, ready <-chan struct{}) *Pruner {
	if batch < 1 {
		batch = 1
	}

	return &Pruner{
		grid:     grid,
		repo:     repo,
		org:      org,
		keep:     keep,
		interval: interval,
		batch:    batch,
		ready:    ready,
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
	}
}

func (p *Pruner) Start() {
	defer close(p.done)

	select {
	case <-p.ready:
	case <-p.stop:
		return
	}

	log.Infof("kpi_windows pruner starting: keeping %s, every %s", p.keep, p.interval)

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		p.PruneOnce(time.Now().UTC())

		select {
		case <-ticker.C:
		case <-p.stop:
			return
		}
	}
}

func (p *Pruner) Stop() {
	close(p.stop)
	<-p.done
}

func (p *Pruner) PruneOnce(now time.Time) {
	before := p.grid.WindowAt(now.Add(-p.keep)).ID

	windows, err := p.drain(func() (int64, error) {
		return p.repo.PruneWindowsBefore(p.org, before, p.batch)
	})
	if err != nil {
		log.Errorf("pruner: kpi_windows before window %d: %v", before, err)

		return
	}

	entries, err := p.drain(func() (int64, error) {
		return p.repo.PruneKpiLedgerBefore(p.org, before, p.batch)
	})
	if err != nil {
		log.Errorf("pruner: kpi ledger before window %d: %v", before, err)

		return
	}

	if windows > 0 || entries > 0 {
		log.Infof("pruner: removed %d kpi_windows rows and %d kpi ledger entries before window %d",
			windows, entries, before)
	}
}

func (p *Pruner) drain(deleteBatch func() (int64, error)) (int64, error) {
	var total int64

	for {
		n, err := deleteBatch()
		total += n

		if err != nil {
			return total, err
		}

		if n < int64(p.batch) {
			return total, nil
		}

		select {
		case <-p.stop:
			return total, nil
		default:
		}
	}
}
