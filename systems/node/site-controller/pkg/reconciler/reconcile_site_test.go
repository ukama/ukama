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
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	uuid "github.com/ukama/ukama/systems/common/uuid"
	contpb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	contmocks "github.com/ukama/ukama/systems/node/controller/pb/gen/mocks"
	scmocks "github.com/ukama/ukama/systems/node/site-controller/mocks"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/adapters"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/policy"
)

const (
	testSiteID  = "44444444-4444-4444-4444-444444444444"
	testTowerID = "uk-sa2633-tnode-v0-0001"
	testAmpID   = "uk-sa2633-anode-v0-0001"
	testCNodeID = "uk-sa2633-cnode-v0-0001"
)

type fakeProvider struct {
	client contpb.ControllerServiceClient
}

func (f *fakeProvider) GetClient() (contpb.ControllerServiceClient, error) {
	return f.client, nil
}

type fakeFlights struct {
	flight  *db.SiteIntentFlight
	upserts []db.SiteIntentFlight
}

func (f *fakeFlights) GetBySiteIntentID(id uuid.UUID) (*db.SiteIntentFlight, error) {
	if f.flight == nil {
		return nil, nil
	}
	current := *f.flight

	return &current, nil
}

func (f *fakeFlights) Upsert(flight *db.SiteIntentFlight) error {
	stored := *flight
	stored.UpdatedAt = time.Now().UTC()
	f.flight = &stored
	f.upserts = append(f.upserts, stored)

	return nil
}

type harness struct {
	reconciler *Reconciler
	intents    *scmocks.IntentRepo
	states     *scmocks.StateRepo
	flights    *fakeFlights
	ports      *scmocks.PortMapRepo
	controller *contmocks.ControllerServiceClient
	intentID   uuid.UUID
}

func newHarness(t *testing.T) *harness {
	t.Helper()

	intents := &scmocks.IntentRepo{}
	states := &scmocks.StateRepo{}
	flights := &fakeFlights{}
	ports := &scmocks.PortMapRepo{}
	controller := &contmocks.ControllerServiceClient{}

	provider := &fakeProvider{client: controller}

	states.On("Upsert", mock.Anything).Return(nil).Maybe()
	controller.On("SendNodeCommand", mock.Anything, mock.Anything).
		Return(&contpb.SendNodeCommandResponse{}, nil).Maybe()

	ports.On("GetBySite", testSiteID).Return([]db.SitePortMap{
		{Port: 1, Role: policy.RoleCNode, NodeID: testCNodeID},
		{Port: 2, Role: policy.RoleTower, NodeID: testTowerID},
		{Port: 3, Role: policy.RoleAmplifier, NodeID: testAmpID},
	}, nil).Maybe()

	r := New(intents, states, flights, ports, nil, provider,
		nil, nil, adapters.NewCNodeAdapter(provider), 30*time.Second, 3)

	return &harness{
		reconciler: r,
		intents:    intents,
		states:     states,
		flights:    flights,
		ports:      ports,
		controller: controller,
		intentID:   uuid.NewV4(),
	}
}

func (h *harness) withIntent(service, radio string) *harness {
	h.intents.On("Get", testSiteID).Return(&db.SiteIntent{
		ID:             h.intentID,
		SiteID:         testSiteID,
		DesiredService: service,
		DesiredRadio:   radio,
	}, nil).Maybe()

	return h
}

func (h *harness) withState(service, radio string) *harness {
	h.states.On("Get", testSiteID).Return(&db.SiteState{
		SiteID:       testSiteID,
		ServiceState: service,
		RadioState:   radio,
	}, nil).Maybe()

	return h
}

func (h *harness) withFlight(status string, retries int, updatedAt time.Time) *harness {
	h.flights.flight = &db.SiteIntentFlight{
		SiteIntentID: h.intentID,
		Status:       status,
		RetryCount:   retries,
		ExpiresAt:    time.Now().UTC().Add(time.Hour),
		UpdatedAt:    updatedAt,
	}

	return h
}

func TestReconcileSite_ReappliesServiceAfterTowerRestart(t *testing.T) {
	h := newHarness(t).
		withIntent(StateOn, StateOff).
		withState(StateOff, StateOff).
		withFlight(db.IntentFlightStatusSucceeded, 0, time.Now().UTC())

	h.controller.On("SendNodeCommand", mock.Anything, mock.Anything).
		Return(&contpb.SendNodeCommandResponse{}, nil).Maybe()
	h.controller.On("ToggleService", mock.Anything, mock.MatchedBy(func(req *contpb.ToggleServiceRequest) bool {
		return req.NodeId == testTowerID && req.State == StateOn
	})).Return(&contpb.ToggleServiceResponse{}, nil).Once()

	_ = h.reconciler.ReconcileSite(context.Background(), testSiteID, false)

	h.controller.AssertExpectations(t)
}

func TestReconcileSite_DoesNotActWhenObservedMatchesDesired(t *testing.T) {
	h := newHarness(t).
		withIntent(StateOn, StateOn).
		withState(StateRunning, StateOn).
		withFlight(db.IntentFlightStatusSucceeded, 0, time.Now().UTC())

	err := h.reconciler.ReconcileSite(context.Background(), testSiteID, false)

	require.NoError(t, err)
	h.controller.AssertNotCalled(t, "ToggleService", mock.Anything, mock.Anything)
	h.controller.AssertNotCalled(t, "ToggleRadio", mock.Anything, mock.Anything)
}

func TestReconcileSite_ReappliesRadioAfterDrift(t *testing.T) {
	h := newHarness(t).
		withIntent(StateOff, StateOn).
		withState(StateOff, StateOff).
		withFlight(db.IntentFlightStatusSucceeded, 0, time.Now().UTC())

	h.controller.On("ToggleRadio", mock.Anything, mock.MatchedBy(func(req *contpb.ToggleRadioRequest) bool {
		return req.NodeId == testAmpID && req.State == StateOn
	})).Return(&contpb.ToggleRadioResponse{}, nil).Once()

	_ = h.reconciler.ReconcileSite(context.Background(), testSiteID, false)

	h.controller.AssertExpectations(t)
}

func TestReconcileSite_ExhaustedFlightBacksOffBeforeRetrying(t *testing.T) {
	h := newHarness(t).
		withIntent(StateOff, StateOn).
		withState(StateOff, StateOff).
		withFlight(db.IntentFlightStatusTimeout, 3, time.Now().UTC())

	err := h.reconciler.ReconcileSite(context.Background(), testSiteID, false)

	require.NoError(t, err)
	h.controller.AssertNotCalled(t, "ToggleRadio", mock.Anything, mock.Anything)
}

func TestReconcileSite_ExhaustedFlightRetriesOnceBackoffElapsed(t *testing.T) {
	h := newHarness(t).
		withIntent(StateOff, StateOn).
		withState(StateOff, StateOff).
		withFlight(db.IntentFlightStatusTimeout, 3, time.Now().UTC().Add(-2*time.Minute))

	h.controller.On("ToggleRadio", mock.Anything, mock.Anything).
		Return(&contpb.ToggleRadioResponse{}, nil).Once()

	_ = h.reconciler.ReconcileSite(context.Background(), testSiteID, false)

	h.controller.AssertExpectations(t)
}

func TestRearmDue(t *testing.T) {
	r := &Reconciler{reconcileInterval: 30 * time.Second}

	assert.True(t, r.rearmDue(nil),
		"no flight means nothing to wait for")

	assert.True(t, r.rearmDue(&db.SiteIntentFlight{
		Status:    db.IntentFlightStatusSucceeded,
		UpdatedAt: time.Now().UTC(),
	}), "drift after a converged flight must be acted on immediately")

	assert.False(t, r.rearmDue(&db.SiteIntentFlight{
		Status:    db.IntentFlightStatusTimeout,
		UpdatedAt: time.Now().UTC(),
	}), "a just-exhausted flight must back off")

	assert.True(t, r.rearmDue(&db.SiteIntentFlight{
		Status:    db.IntentFlightStatusTimeout,
		UpdatedAt: time.Now().UTC().Add(-time.Minute),
	}), "an exhausted flight must retry once the interval has passed")
}

func TestMarkConverged_ResetsRetryCount(t *testing.T) {
	flights := &fakeFlights{}
	intentID := uuid.NewV4()

	r := &Reconciler{flights: flights, flightTTL: time.Hour}

	err := r.markConverged(&db.SiteIntent{ID: intentID}, &db.SiteIntentFlight{
		SiteIntentID: intentID,
		Status:       db.IntentFlightStatusTimeout,
		RetryCount:   3,
	})

	require.NoError(t, err)
	require.Len(t, flights.upserts, 1)
	assert.Equal(t, db.IntentFlightStatusSucceeded, flights.upserts[0].Status)
	assert.Equal(t, 0, flights.upserts[0].RetryCount)
}

func TestMarkConverged_SkipsWriteWhenAlreadyConverged(t *testing.T) {
	flights := &fakeFlights{}
	intentID := uuid.NewV4()

	r := &Reconciler{flights: flights, flightTTL: time.Hour}

	err := r.markConverged(&db.SiteIntent{ID: intentID}, &db.SiteIntentFlight{
		SiteIntentID: intentID,
		Status:       db.IntentFlightStatusSucceeded,
		RetryCount:   0,
	})

	require.NoError(t, err)
	assert.Empty(t, flights.upserts, "a converged flight must not be rewritten")
}

func TestSyncStatus(t *testing.T) {
	inSyncIntent := &db.SiteIntent{DesiredService: StateOn, DesiredRadio: StateOn}
	inSyncState := &db.SiteState{ServiceState: StateRunning, RadioState: StateOn}
	driftState := &db.SiteState{ServiceState: StateOff, RadioState: StateOn}

	t.Run("observed matching desired is in sync", func(t *testing.T) {
		status, reason := syncStatus(inSyncIntent, inSyncState, &db.SiteIntentFlight{
			Status: db.IntentFlightStatusSucceeded,
		})

		assert.Equal(t, SyncStatusInSync, status)
		assert.Equal(t, ReasonOK, reason)
	})

	t.Run("in sync wins even after a failed flight", func(t *testing.T) {
		status, _ := syncStatus(inSyncIntent, inSyncState, &db.SiteIntentFlight{
			Status:     db.IntentFlightStatusTimeout,
			RetryCount: 3,
		})

		assert.Equal(t, SyncStatusInSync, status,
			"fault must clear as soon as observed matches desired")
	})

	t.Run("drift while still retrying is applying", func(t *testing.T) {
		status, reason := syncStatus(inSyncIntent, driftState, &db.SiteIntentFlight{
			Status:     db.IntentFlightStatusPending,
			RetryCount: 1,
		})

		assert.Equal(t, SyncStatusApplying, status)
		assert.Contains(t, reason, "service desired=on observed=off")
	})

	t.Run("drift after retries are exhausted is a fault", func(t *testing.T) {
		for _, terminal := range []string{
			db.IntentFlightStatusFailed,
			db.IntentFlightStatusTimeout,
			db.IntentFlightStatusExpired,
		} {
			status, reason := syncStatus(inSyncIntent, driftState, &db.SiteIntentFlight{
				Status:     terminal,
				RetryCount: 3,
			})

			assert.Equal(t, SyncStatusFault, status, "%s must surface as fault", terminal)
			assert.Contains(t, reason, "attempts=3")
		}
	})

	t.Run("drift with no flight yet is applying", func(t *testing.T) {
		status, _ := syncStatus(inSyncIntent, driftState, nil)

		assert.Equal(t, SyncStatusApplying, status)
	})

	t.Run("radio drift is named in the reason", func(t *testing.T) {
		_, reason := syncStatus(
			&db.SiteIntent{DesiredService: StateOff, DesiredRadio: StateOn},
			&db.SiteState{ServiceState: StateOff, RadioState: StateOff},
			&db.SiteIntentFlight{Status: db.IntentFlightStatusTimeout, RetryCount: 2},
		)

		assert.Contains(t, reason, "radio desired=on observed=off")
	})
}

func TestApplyIntentState_ServiceAppliedWhenRadioFails(t *testing.T) {
	h := newHarness(t).
		withIntent(StateOn, StateOn).
		withState(StateOff, StateOff).
		withFlight(db.IntentFlightStatusSucceeded, 0, time.Now().UTC())

	h.controller.On("SendNodeCommand", mock.Anything, mock.Anything).
		Return(&contpb.SendNodeCommandResponse{}, nil).Maybe()
	h.controller.On("ToggleRadio", mock.Anything, mock.Anything).
		Return(nil, assert.AnError).Once()
	h.controller.On("ToggleService", mock.Anything, mock.MatchedBy(func(req *contpb.ToggleServiceRequest) bool {
		return req.NodeId == testTowerID && req.State == StateOn
	})).Return(&contpb.ToggleServiceResponse{}, nil).Once()

	_ = h.reconciler.ReconcileSite(context.Background(), testSiteID, true)

	h.controller.AssertExpectations(t)
}

func poePortsFor(controller *contmocks.ControllerServiceClient) []string {
	var paths []string
	for _, c := range controller.Calls {
		if c.Method != "SendNodeCommand" {
			continue
		}
		req, ok := c.Arguments.Get(1).(*contpb.SendNodeCommandRequest)
		if !ok {
			continue
		}
		if strings.Contains(req.Path, "/poe") {
			paths = append(paths, req.Path)
		}
	}
	return paths
}

func TestServiceOnDoesNotPowerAmplifierPort(t *testing.T) {
	h := newHarness(t).
		withIntent(StateOn, StateOff).
		withState(StateOff, StateOff).
		withFlight(db.IntentFlightStatusSucceeded, 0, time.Now().UTC())

	h.controller.On("ToggleService", mock.Anything, mock.Anything).
		Return(&contpb.ToggleServiceResponse{}, nil).Maybe()

	_ = h.reconciler.ReconcileSite(context.Background(), testSiteID, true)

	paths := poePortsFor(h.controller)
	assert.Contains(t, paths, "/v1/ports/2/poe", "tower port must be powered for service on")
	assert.NotContains(t, paths, "/v1/ports/3/poe", "amplifier port belongs to radio intent, not service")
}

func TestRadioOnPowersAmplifierPort(t *testing.T) {
	h := newHarness(t).
		withIntent(StateOff, StateOn).
		withState(StateOff, StateOff).
		withFlight(db.IntentFlightStatusSucceeded, 0, time.Now().UTC())

	h.controller.On("ToggleRadio", mock.Anything, mock.Anything).
		Return(&contpb.ToggleRadioResponse{}, nil).Maybe()

	_ = h.reconciler.ReconcileSite(context.Background(), testSiteID, true)

	assert.Contains(t, poePortsFor(h.controller), "/v1/ports/3/poe",
		"amplifier port must be powered for radio on")
}
