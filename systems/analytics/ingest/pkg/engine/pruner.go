/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package engine

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/analytics/ingest/pkg/db"
	"github.com/ukama/ukama/systems/analytics/schema"
)

type Pruner struct {
	grid     schema.Grid
	datasets []string
	raw      db.RawRepo
	ledger   db.LedgerRepo
	org      string
	keep     int64
	interval time.Duration
	batch    int
	stop     chan struct{}
	done     chan struct{}
}

func NewPruner(grid schema.Grid, specs []schema.SourceSpec, raw db.RawRepo,
	ledger db.LedgerRepo, org string, keep int64, interval time.Duration, batch int) *Pruner {
	datasets := []string{}

	for _, s := range specs {
		for _, p := range s.Pulls {
			if p.Strategy == schema.StrategyWindow {
				datasets = append(datasets, p.Key)
			}
		}
	}

	if batch < 1 {
		batch = 1
	}

	return &Pruner{
		grid:     grid,
		datasets: datasets,
		raw:      raw,
		ledger:   ledger,
		org:      org,
		keep:     keep,
		interval: interval,
		batch:    batch,
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
	}
}

func (p *Pruner) Datasets() []string {
	return p.datasets
}

func (p *Pruner) Start() {
	defer close(p.done)

	log.Infof("pruner starting: keeping %d windows of %v, every %s", p.keep, p.datasets, p.interval)

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
	before := p.grid.NewestEligible(now) - p.keep + 1

	for _, key := range p.datasets {
		rows, err := p.drain(func() (int64, error) {
			return p.raw.PruneBefore(p.org, key, before, p.batch)
		})
		if err != nil {
			log.Errorf("pruner: dataset %s rows before window %d: %v", key, before, err)

			continue
		}

		entries, err := p.drain(func() (int64, error) {
			return p.ledger.PruneBefore(p.org, schema.LedgerKindDataset, key, before, p.batch)
		})
		if err != nil {
			log.Errorf("pruner: dataset %s ledger before window %d: %v", key, before, err)

			continue
		}

		if rows > 0 || entries > 0 {
			log.Infof("pruner: dataset %s removed %d rows and %d ledger entries before window %d",
				key, rows, entries, before)
		}
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
