/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
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
	repo     db.OperationRepo
	interval time.Duration
	batch    int
}

func NewSweeper(repo db.OperationRepo) *Sweeper {
	return &Sweeper{
		repo:     repo,
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
			s.sweepOnce()
		}
	}
}

func (s *Sweeper) sweepOnce() {
	expired, err := s.repo.FindExpired(time.Now().UTC(), s.batch)
	if err != nil {
		log.Errorf("sweeper: FindExpired error: %v", err)
		return
	}
	if len(expired) == 0 {
		return
	}
	log.Infof("sweeper: timing out %d expired operation(s)", len(expired))
	for i := range expired {
		op := &expired[i]
		_, err := s.repo.Terminate(op.Id, op.FencingToken, db.OperationTimeout, db.OperationAudit{
			Event:  "timeout",
			Reason: "lease expired",
		}, "lease expired")
		if err != nil {
			log.Warnf("sweeper: terminate(timeout) for %s failed: %v", op.Id, err)
		}
	}
}
