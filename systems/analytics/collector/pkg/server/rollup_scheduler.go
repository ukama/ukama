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

	"github.com/ukama/ukama/systems/analytics/collector/pkg/db"

	log "github.com/sirupsen/logrus"
)

type RollupSchedulerConfig struct {
	Enabled      bool
	Interval     time.Duration
	LookbackDays int
}

type RollupScheduler struct {
	stateRepo  db.StateRepo
	rollupRepo db.RollupRepo
	config     RollupSchedulerConfig
}

func NewRollupScheduler(stateRepo db.StateRepo, rollupRepo db.RollupRepo,
	config RollupSchedulerConfig) *RollupScheduler {
	if config.Interval <= 0 {
		config.Interval = 5 * time.Minute
	}
	if config.LookbackDays <= 0 {
		config.LookbackDays = 30
	}

	return &RollupScheduler{
		stateRepo:  stateRepo,
		rollupRepo: rollupRepo,
		config:     config,
	}
}

func (s *RollupScheduler) Start(ctx context.Context) {
	if !s.config.Enabled {
		log.Info("rollup scheduler disabled")
		return
	}

	log.Infof("rollup scheduler enabled: interval=%s lookback_days=%d",
		s.config.Interval, s.config.LookbackDays)

	go func() {
		s.rebuildDirty()

		ticker := time.NewTicker(s.config.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Info("rollup scheduler stopped")
				return
			case <-ticker.C:
				s.rebuildDirty()
			}
		}
	}()
}

func (s *RollupScheduler) rebuildDirty() {
	states, err := s.stateRepo.GetRollupStates()
	if err != nil {
		log.Errorf("failed to read rollup states: %v", err)
		return
	}

	to := time.Now().UTC()
	from := to.AddDate(0, 0, -s.config.LookbackDays)

	for i := range states {
		if !states[i].Dirty {
			continue
		}

		if err := s.rebuild(states[i].Rollup, from, to); err != nil {
			log.Errorf("failed to rebuild dirty rollup %s: %v",
				states[i].Rollup, err)
			continue
		}

		if err := s.stateRepo.SetRollupWatermark(states[i].Rollup, to); err != nil {
			log.Errorf("failed to update watermark for rollup %s: %v",
				states[i].Rollup, err)
		}
	}
}

func (s *RollupScheduler) rebuild(name string, from, to time.Time) error {
	switch name {
	case "business_sales_daily":
		return s.rollupRepo.RebuildSalesDaily(from, to)
	case "business_package_daily":
		return s.rollupRepo.RebuildPackageDaily(from, to)
	case "business_billing_daily":
		return s.rollupRepo.RebuildBillingDaily(from, to)
	case "customer_usage_daily":
		return s.rollupRepo.RebuildCustomerUsageDaily(from, to)
	case "customer_state_daily":
		return s.rollupRepo.RebuildCustomerStateDaily(from, to)
	case "alarm_daily":
		return s.rollupRepo.RebuildAlarmDaily(from, to)
	case "metric_hourly":
		return s.rollupRepo.RebuildMetricHourly(from, to)
	default:
		// TODO(analytics-phase-2): implement rebuild support for
		// business_site_daily, business_inventory_daily,
		// network_health_hourly, site_health_hourly, node_health_hourly,
		// radio_hourly, backhaul_hourly and power_hourly.
		log.Warnf("rollup %s is dirty but has no rebuild implementation", name)
		return nil
	}
}
