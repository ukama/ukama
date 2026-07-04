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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/analytics/collector/pkg"
	"github.com/ukama/ukama/systems/analytics/collector/pkg/clients"
	"github.com/ukama/ukama/systems/analytics/collector/pkg/db"
	"github.com/ukama/ukama/systems/analytics/collector/pkg/refresh"
	"github.com/ukama/ukama/systems/common/uuid"

	pb "github.com/ukama/ukama/systems/analytics/collector/pb/gen"
)

const testOrgName = "testorg"

/* Hand-written stub repos (mockery mocks are generated later by make gen). */

type stubStateRepo struct {
	refreshStates map[string]db.RefreshState
	rollupStates  map[string]db.RollupState
}

func newStubStateRepo() *stubStateRepo {
	return &stubStateRepo{
		refreshStates: map[string]db.RefreshState{},
		rollupStates:  map[string]db.RollupState{},
	}
}

func (s *stubStateRepo) UpsertRefreshState(state *db.RefreshState) error {
	s.refreshStates[state.Source] = *state

	return nil
}

func (s *stubStateRepo) GetRefreshStates() ([]db.RefreshState, error) {
	out := make([]db.RefreshState, 0, len(s.refreshStates))
	for _, v := range s.refreshStates {
		out = append(out, v)
	}

	return out, nil
}

func (s *stubStateRepo) MarkRollupDirty(rollup string) error {
	st := s.rollupStates[rollup]
	st.Rollup = rollup
	st.Dirty = true
	s.rollupStates[rollup] = st

	return nil
}

func (s *stubStateRepo) SetRollupWatermark(rollup string, watermark time.Time) error {
	st := s.rollupStates[rollup]
	st.Rollup = rollup
	st.Watermark = watermark
	st.Dirty = false
	s.rollupStates[rollup] = st

	return nil
}

func (s *stubStateRepo) GetRollupStates() ([]db.RollupState, error) {
	out := make([]db.RollupState, 0, len(s.rollupStates))
	for _, v := range s.rollupStates {
		out = append(out, v)
	}

	return out, nil
}

type stubSnapshotRepo struct {
	customers map[string]db.CustomerSnapshot
	sims      map[string]db.SimSnapshot
	packages  map[string]db.PackageSnapshot
}

func newStubSnapshotRepo() *stubSnapshotRepo {
	return &stubSnapshotRepo{
		customers: map[string]db.CustomerSnapshot{},
		sims:      map[string]db.SimSnapshot{},
		packages:  map[string]db.PackageSnapshot{},
	}
}

func (s *stubSnapshotRepo) UpsertNetwork(v *db.NetworkSnapshot) error { return nil }
func (s *stubSnapshotRepo) UpsertSite(v *db.SiteSnapshot) error       { return nil }
func (s *stubSnapshotRepo) UpsertNode(v *db.NodeSnapshot) error       { return nil }
func (s *stubSnapshotRepo) UpdateNodeStatus(nodeId, nodeStatus string, at time.Time) error {
	return nil
}
func (s *stubSnapshotRepo) UpsertCustomer(v *db.CustomerSnapshot) error {
	s.customers[v.CustomerId.String()] = *v

	return nil
}
func (s *stubSnapshotRepo) DeleteCustomer(customerId string) error {
	delete(s.customers, customerId)

	return nil
}
func (s *stubSnapshotRepo) UpsertSim(v *db.SimSnapshot) error {
	s.sims[v.SimId] = *v

	return nil
}
func (s *stubSnapshotRepo) DeleteSim(simId string) error {
	delete(s.sims, simId)

	return nil
}
func (s *stubSnapshotRepo) UpsertSimBatch(v *db.SimBatchSnapshot) error { return nil }
func (s *stubSnapshotRepo) UpsertPackage(v *db.PackageSnapshot) error {
	s.packages[v.PackageId.String()] = *v

	return nil
}
func (s *stubSnapshotRepo) DeletePackage(packageId string) error {
	delete(s.packages, packageId)

	return nil
}
func (s *stubSnapshotRepo) UpsertInventory(v *db.InventorySnapshot) error       { return nil }
func (s *stubSnapshotRepo) UpsertBilling(v *db.BillingSnapshot) error           { return nil }
func (s *stubSnapshotRepo) UpsertHealthReport(v *db.HealthReportSnapshot) error { return nil }

type stubFactRepo struct {
	payments []db.PaymentEvent
}

func (s *stubFactRepo) AddPaymentEvent(e *db.PaymentEvent) error {
	for _, p := range s.payments {
		if p.ExternalId == e.ExternalId {
			/* mirror ON CONFLICT DO NOTHING */
			return nil
		}
	}

	s.payments = append(s.payments, *e)

	return nil
}
func (s *stubFactRepo) AddUsageEvent(e *db.UsageEvent) error         { return nil }
func (s *stubFactRepo) AddMetricSample(e *db.MetricSample) error     { return nil }
func (s *stubFactRepo) AddAlarmEvent(e *db.AlarmEvent) error         { return nil }
func (s *stubFactRepo) AddNodeStateEvent(e *db.NodeStateEvent) error { return nil }
func (s *stubFactRepo) AddSiteStateEvent(e *db.SiteStateEvent) error { return nil }
func (s *stubFactRepo) AddCustomerEvent(e *db.CustomerEvent) error   { return nil }
func (s *stubFactRepo) AddSimEvent(e *db.SimEvent) error             { return nil }
func (s *stubFactRepo) AddPackageEvent(e *db.PackageEvent) error     { return nil }
func (s *stubFactRepo) AddInventoryEvent(e *db.InventoryEvent) error { return nil }
func (s *stubFactRepo) TransitionNodeState(nodeId, state string, at time.Time) error {
	return nil
}
func (s *stubFactRepo) TransitionSiteState(siteId uuid.UUID, state string, at time.Time) error {
	return nil
}
func (s *stubFactRepo) TransitionSimState(simId, state string, at time.Time) error {
	return nil
}
func (s *stubFactRepo) OpenCustomerPackageInterval(customerId, packageId uuid.UUID, state string, at time.Time) error {
	return nil
}
func (s *stubFactRepo) CloseCustomerPackageInterval(customerId uuid.UUID, at time.Time) error {
	return nil
}

type stubEventRepo struct {
	seen   map[string]bool
	errors []db.EventError
}

func newStubEventRepo() *stubEventRepo {
	return &stubEventRepo{seen: map[string]bool{}}
}

func (s *stubEventRepo) LogEvent(l *db.EventLog) (bool, error) {
	if s.seen[l.MsgId] {
		return false, nil
	}

	s.seen[l.MsgId] = true

	return true, nil
}

func (s *stubEventRepo) RecordError(e *db.EventError) error {
	s.errors = append(s.errors, *e)

	return nil
}

func (s *stubEventRepo) GetRecent(limit int) ([]db.EventLog, error) {
	return nil, nil
}

type stubRegistryClient struct{ fail bool }

func (c *stubRegistryClient) GetNetworks() ([]clients.RegistryNetwork, error) {
	if c.fail {
		return nil, assert.AnError
	}

	return []clients.RegistryNetwork{}, nil
}

func (c *stubRegistryClient) GetSites() ([]clients.RegistrySite, error) {
	if c.fail {
		return nil, assert.AnError
	}

	return []clients.RegistrySite{}, nil
}

func newTestRefresher(stateRepo db.StateRepo, registryClient clients.RegistryClient) *refresh.Refresher {
	return refresh.NewRefresher(stateRepo, newStubSnapshotRepo(), &stubFactRepo{},
		nil, nil, nil, nil, nil, nil, "")
}

func TestCollectorServer_Refresh(t *testing.T) {
	t.Run("RefreshRegistrySucceeds", func(t *testing.T) {
		stateRepo := newStubStateRepo()
		s := NewCollectorServer(testOrgName, stateRepo, nil, newStubEventRepo(),
			newTestRefresher(stateRepo, &stubRegistryClient{}), nil, "")

		resp, err := s.Refresh(context.TODO(), &pb.RefreshRequest{Source: "registry"})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.States, 1)
		assert.Equal(t, "registry", resp.States[0].Source)
		assert.Equal(t, "ok", resp.States[0].Status)

		state := stateRepo.refreshStates["registry"]
		assert.Equal(t, "ok", state.Status)
		assert.False(t, state.LastRunAt.IsZero())
		assert.False(t, state.LastSuccessAt.IsZero())
	})

	t.Run("RefreshRegistryFails", func(t *testing.T) {
		stateRepo := newStubStateRepo()
		s := NewCollectorServer(testOrgName, stateRepo, nil, newStubEventRepo(),
			newTestRefresher(stateRepo, &stubRegistryClient{fail: true}), nil, "")

		resp, err := s.Refresh(context.TODO(), &pb.RefreshRequest{Source: "registry"})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.States, 1)
		assert.Equal(t, "failed", resp.States[0].Status)
		assert.NotEmpty(t, resp.States[0].Detail)

		state := stateRepo.refreshStates["registry"]
		assert.Equal(t, "failed", state.Status)
		assert.True(t, state.LastSuccessAt.IsZero())
	})

	t.Run("RefreshInvalidSource", func(t *testing.T) {
		stateRepo := newStubStateRepo()
		s := NewCollectorServer(testOrgName, stateRepo, nil, newStubEventRepo(),
			newTestRefresher(stateRepo, &stubRegistryClient{}), nil, "")

		resp, err := s.Refresh(context.TODO(), &pb.RefreshRequest{Source: "bogus"})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})
}

func TestCollectorServer_SeedDemo(t *testing.T) {
	t.Run("DeniedWhenNotDebug", func(t *testing.T) {
		pkg.IsDebugMode = false

		stateRepo := newStubStateRepo()
		s := NewCollectorServer(testOrgName, stateRepo, nil, newStubEventRepo(),
			newTestRefresher(stateRepo, &stubRegistryClient{}), nil, "")

		resp, err := s.SeedDemo(context.TODO(), &pb.SeedDemoRequest{Sites: 1})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.PermissionDenied, status.Code(err))
	})

	t.Run("AllowedInDebug", func(t *testing.T) {
		pkg.IsDebugMode = true
		defer func() { pkg.IsDebugMode = false }()

		stateRepo := newStubStateRepo()
		s := NewCollectorServer(testOrgName, stateRepo, nil, newStubEventRepo(),
			newTestRefresher(stateRepo, &stubRegistryClient{}), nil, "")

		resp, err := s.SeedDemo(context.TODO(), &pb.SeedDemoRequest{Sites: 1})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, stateRepo.rollupStates["business_sales_daily"].Dirty)
	})
}
