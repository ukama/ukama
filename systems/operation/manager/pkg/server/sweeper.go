/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/operation/manager/pkg"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
)

type Sweeper struct {
	repos    Repos
	interval time.Duration
	batch    int
}

func NewSweeper(repos Repos) *Sweeper {
	return &Sweeper{
		repos:    repos,
		interval: pkg.SweeperInterval,
		batch:    100,
	}
}

func (s *Sweeper) Run(ctx context.Context) {
	log.Infof("operation/manager sweeper started (interval=%s)", s.interval)
	t := time.NewTicker(s.interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Info("operation/manager sweeper stopped")
			return
		case <-t.C:
			s.sweepOnce(ctx)
		}
	}
}

func (s *Sweeper) sweepOnce(parent context.Context) {
	ctx, cancel := context.WithTimeout(parent, s.interval)
	defer cancel()

	now := time.Now().UTC()
	expired, err := s.repos.Operations.List(ctx, db.OperationFilter{
		Statuses:           []db.OperationStatus{db.OperationPending, db.OperationRunning},
		LeaseExpiresBefore: &now,
		OrderBy:            db.OrderByLeaseExpiry,
		Limit:              s.batch,
	})
	if err != nil {
		log.Errorf("sweeper: list expired operations: %v", err)
		return
	}
	if len(expired) == 0 {
		return
	}
	log.Infof("sweeper: timing out %d expired operation(s)", len(expired))
	for _, op := range expired {
		_, err := s.repos.terminate(ctx, op.Id, op.FencingToken, db.OperationTimeout,
			db.OperationAudit{
				Event:  "timeout",
				Actor:  pkg.OperationManagerSweeperActor,
				Reason: "lease expired",
			}, "lease expired")
		if err != nil {
			log.Warnf("sweeper: terminate(timeout) for %s failed: %v", op.Id, err)
		}
	}
}
