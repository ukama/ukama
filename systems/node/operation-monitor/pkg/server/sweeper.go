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
	"google.golang.org/protobuf/types/known/timestamppb"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"

	"github.com/ukama/ukama/systems/node/operation-monitor/pkg"
	"github.com/ukama/ukama/systems/node/operation-monitor/pkg/db"
)

type Sweeper struct {
	monitor  *MonitorServer
	interval time.Duration
	batch    int
}

func NewSweeper(monitor *MonitorServer) *Sweeper {
	return &Sweeper{
		monitor:  monitor,
		interval: pkg.SweeperInterval,
		batch:    100,
	}
}

func (s *Sweeper) Run(ctx context.Context) {
	log.Infof("operation-monitor sweeper started (interval=%s)", s.interval)
	t := time.NewTicker(s.interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Info("operation-monitor sweeper stopped")
			return
		case <-t.C:
			s.sweepOnce()
		}
	}
}

func (s *Sweeper) sweepOnce() {
	expired, err := s.monitor.repo.FindExpired(time.Now().UTC(), s.batch)
	if err != nil {
		log.Errorf("sweeper: FindExpired error: %v", err)
		return
	}
	if len(expired) == 0 {
		return
	}
	log.Infof("sweeper: expiring %d intent(s) past deadline", len(expired))
	for i := range expired {
		intent := &expired[i]
		if _, err := s.monitor.repo.MarkTerminal(intent.OperationId, db.IntentExpired); err != nil {
			log.Warnf("sweeper: mark expired for %s failed: %v", intent.OperationId, err)
			continue
		}
		if err := s.publishFailed(intent, "deadline exceeded"); err != nil {
			log.Warnf("sweeper: publish failed for %s: %v", intent.OperationId, err)
		}
	}
}

func (s *Sweeper) publishFailed(intent *db.MonitoredIntent, reason string) error {
	route := s.monitor.publishBuilder.SetAction("failed").SetObject("operation").MustBuild()
	return s.monitor.msgbus.PublishRequest(route, &epb.OperationFailedEvent{
		OperationId:  intent.OperationId.String(),
		FencingToken: intent.FencingToken,
		ResourceKey:  intent.ResourceKey,
		Reason:       reason,
		FailedAt:     timestamppb.Now(),
	})
}
