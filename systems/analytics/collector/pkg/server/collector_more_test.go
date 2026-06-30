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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"

	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"

	pb "github.com/ukama/ukama/systems/analytics/collector/pb/gen"
	"github.com/ukama/ukama/systems/analytics/collector/pkg/db"
)

type stubRollupRepo struct{ calls int }

func (r *stubRollupRepo) RebuildSalesDaily(from, to time.Time) error         { r.calls++; return nil }
func (r *stubRollupRepo) RebuildPackageDaily(from, to time.Time) error       { r.calls++; return nil }
func (r *stubRollupRepo) RebuildBillingDaily(from, to time.Time) error       { r.calls++; return nil }
func (r *stubRollupRepo) RebuildCustomerUsageDaily(from, to time.Time) error { r.calls++; return nil }
func (r *stubRollupRepo) RebuildCustomerStateDaily(from, to time.Time) error { r.calls++; return nil }
func (r *stubRollupRepo) RebuildAlarmDaily(from, to time.Time) error         { r.calls++; return nil }
func (r *stubRollupRepo) RebuildMetricHourly(from, to time.Time) error       { r.calls++; return nil }

func (r *stubRollupRepo) UpsertBusinessSalesDaily(*db.BusinessSalesRollupDaily) error         { return nil }
func (r *stubRollupRepo) UpsertBusinessPackageDaily(*db.BusinessPackageRollupDaily) error     { return nil }
func (r *stubRollupRepo) UpsertBusinessSiteDaily(*db.BusinessSiteRollupDaily) error           { return nil }
func (r *stubRollupRepo) UpsertBusinessInventoryDaily(*db.BusinessInventoryRollupDaily) error { return nil }
func (r *stubRollupRepo) UpsertBusinessBillingDaily(*db.BusinessBillingRollupDaily) error     { return nil }
func (r *stubRollupRepo) UpsertCustomerUsageDaily(*db.CustomerUsageRollupDaily) error         { return nil }
func (r *stubRollupRepo) UpsertCustomerStateDaily(*db.CustomerStateRollupDaily) error         { return nil }
func (r *stubRollupRepo) UpsertNetworkHealthHourly(*db.NetworkHealthRollupHourly) error       { return nil }
func (r *stubRollupRepo) UpsertSiteHealthHourly(*db.SiteHealthRollupHourly) error             { return nil }
func (r *stubRollupRepo) UpsertNodeHealthHourly(*db.NodeHealthRollupHourly) error             { return nil }
func (r *stubRollupRepo) UpsertMetricHourly(*db.MetricRollupHourly) error                     { return nil }
func (r *stubRollupRepo) UpsertAlarmDaily(*db.AlarmRollupDaily) error                         { return nil }
func (r *stubRollupRepo) UpsertRadioHourly(*db.RadioRollupHourly) error                       { return nil }
func (r *stubRollupRepo) UpsertBackhaulHourly(*db.BackhaulRollupHourly) error                 { return nil }
func (r *stubRollupRepo) UpsertPowerHourly(*db.PowerRollupHourly) error                       { return nil }

func TestCollectorServer_GetRefreshState(t *testing.T) {
	stateRepo := newStubStateRepo()
	s := NewCollectorServer(testOrgName, stateRepo, nil, newStubEventRepo(), nil, nil, "")

	resp, err := s.GetRefreshState(context.TODO(), &pb.GetRefreshStateRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestCollectorServer_RebuildRollups(t *testing.T) {
	t.Run("All", func(t *testing.T) {
		stateRepo := newStubStateRepo()
		rollup := &stubRollupRepo{}
		s := NewCollectorServer(testOrgName, stateRepo, rollup, newStubEventRepo(), nil, nil, "")

		resp, err := s.RebuildRollups(context.TODO(), &pb.RebuildRollupsRequest{Family: familyAll})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 7, rollup.calls)
	})

	t.Run("InvalidFamily", func(t *testing.T) {
		s := NewCollectorServer(testOrgName, newStubStateRepo(), &stubRollupRepo{}, newStubEventRepo(), nil, nil, "")

		_, err := s.RebuildRollups(context.TODO(), &pb.RebuildRollupsRequest{Family: "bogus"})
		assert.Error(t, err)
	})
}

func TestRollupScheduler_New_Defaults(t *testing.T) {
	s := NewRollupScheduler(newStubStateRepo(), &stubRollupRepo{}, RollupSchedulerConfig{})

	assert.Equal(t, 5*time.Minute, s.config.Interval)
	assert.Equal(t, 30, s.config.LookbackDays)
}

func TestRollupScheduler_Rebuild(t *testing.T) {
	rollup := &stubRollupRepo{}
	s := NewRollupScheduler(newStubStateRepo(), rollup, RollupSchedulerConfig{})

	now := time.Now()
	for _, name := range []string{
		"business_sales_daily", "business_package_daily", "business_billing_daily",
		"customer_usage_daily", "customer_state_daily", "alarm_daily", "metric_hourly",
	} {
		assert.NoError(t, s.rebuild(name, now.AddDate(0, 0, -1), now))
	}

	assert.Error(t, s.rebuild("unknown_rollup", now, now))
}

func TestRollupScheduler_RebuildDirty(t *testing.T) {
	stateRepo := newStubStateRepo()
	assert.NoError(t, stateRepo.MarkRollupDirty("business_sales_daily"))

	rollup := &stubRollupRepo{}
	s := NewRollupScheduler(stateRepo, rollup, RollupSchedulerConfig{})

	s.rebuildDirty()

	assert.GreaterOrEqual(t, rollup.calls, 1)
}

func TestRollupScheduler_Start(t *testing.T) {
	t.Run("Disabled", func(t *testing.T) {
		s := NewRollupScheduler(newStubStateRepo(), &stubRollupRepo{}, RollupSchedulerConfig{Enabled: false})
		s.Start(context.Background())
	})

	t.Run("Enabled", func(t *testing.T) {
		stateRepo := newStubStateRepo()
		s := NewRollupScheduler(stateRepo, &stubRollupRepo{},
			RollupSchedulerConfig{Enabled: true, Interval: 10 * time.Millisecond, LookbackDays: 7})

		ctx, cancel := context.WithCancel(context.Background())
		s.Start(ctx)
		time.Sleep(25 * time.Millisecond)
		cancel()
		time.Sleep(10 * time.Millisecond)
	})
}

// txEventRepo implements eventTransactionRunner so EventNotification exercises
// its transactional path.
type txEventRepo struct {
	*stubEventRepo
}

func (r *txEventRepo) InTransaction(fn func(db.EventRepo, db.StateRepo, db.SnapshotRepo, db.FactRepo) error) error {
	return fn(r.stubEventRepo, newStubStateRepo(), newStubSnapshotRepo(), &stubFactRepo{})
}

// failingTxRepo runs a transaction whose inner dispatch fails (LogEvent errors),
// so EventNotification exercises recordProcessingError.
type failingTxRepo struct {
	*stubEventRepo
}

func (r *failingTxRepo) LogEvent(*db.EventLog) (bool, error) {
	return false, errors.New("log boom")
}

func (r *failingTxRepo) InTransaction(fn func(db.EventRepo, db.StateRepo, db.SnapshotRepo, db.FactRepo) error) error {
	return fn(r, newStubStateRepo(), newStubSnapshotRepo(), &stubFactRepo{})
}

func TestEventNotification_TransactionError(t *testing.T) {
	s := NewCollectorEventServer(testOrgName, &failingTxRepo{newStubEventRepo()},
		newStubStateRepo(), newStubSnapshotRepo(), &stubFactRepo{})

	anyMsg, _ := anypb.New(&epb.Payment{Id: "p-err", Status: "success", AmountCents: 1, ItemType: "package"})

	_, err := s.EventNotification(context.TODO(), &epb.Event{
		RoutingKey: msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventPaymentSuccess]),
		Msg:        anyMsg,
	})

	assert.Error(t, err)
}

func TestEventNotification_TransactionPath(t *testing.T) {
	s := NewCollectorEventServer(testOrgName, &txEventRepo{newStubEventRepo()},
		newStubStateRepo(), newStubSnapshotRepo(), &stubFactRepo{})

	payment := &epb.Payment{Id: "pay-tx", Status: "success", AmountCents: 100, ItemType: "package"}
	anyMsg, _ := anypb.New(payment)

	resp, err := s.EventNotification(context.TODO(), &epb.Event{
		RoutingKey: msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventPaymentSuccess]),
		Msg:        anyMsg,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
