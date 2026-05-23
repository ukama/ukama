/*
* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at https://mozilla.org/MPL/2.0/.
*
* Copyright (c) 2026-present, Ukama Inc.
 */

package reconciler

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
)

type Worker struct {
	reconciler *Reconciler
	sites      db.SiteRepo
	interval   time.Duration
}

func NewWorker(r *Reconciler, sites db.SiteRepo, interval time.Duration) *Worker {
	return &Worker{reconciler: r, sites: sites, interval: interval}
}

func (w *Worker) Start(ctx context.Context) {
	if w.interval <= 0 {
		log.Warn("site-controller: reconcile worker disabled (interval <= 0)")
		return
	}
	log.Infof("site-controller: starting intent reconcile worker (interval=%s)", w.interval)
	go w.run(ctx)
}

func (w *Worker) run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.tick(ctx)
	for {
		select {
		case <-ctx.Done():
			log.Info("site-controller: stopping intent reconcile worker")
			return
		case <-ticker.C:
			w.tick(ctx)
		}
	}
}

func (w *Worker) tick(ctx context.Context) {
	sites, err := w.sites.List()
	if err != nil {
		log.Errorf("site-controller: list sites for reconcile: %v", err)
		return
	}
	for _, site := range sites {
		if err := w.reconciler.ReconcileSite(ctx, site.SiteID, false); err != nil {
			log.Warnf("site-controller: reconcile site %s: %v", site.SiteID, err)
		}
	}
}
