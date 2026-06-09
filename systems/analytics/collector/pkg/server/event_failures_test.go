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

	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"

	"github.com/ukama/ukama/systems/analytics/collector/mocks"
)

// failingServer builds an event server where exactly one repo layer fails,
// so handlers exercise their recordFailure branch at that layer.
func failingServer(layer string) *CollectorEventServer {
	boom := errors.New("boom")
	errFor := func(l string) error {
		if l == layer {
			return boom
		}
		return nil
	}

	ev := &mocks.EventRepo{}
	ev.On("LogEvent", mock.Anything).Return(true, nil).Maybe()
	ev.On("RecordError", mock.Anything).Return(nil).Maybe()

	snap := &mocks.SnapshotRepo{}
	for _, m := range []string{
		"UpsertNetwork", "UpsertSite", "UpsertNode", "UpsertCustomer", "UpsertSim",
		"UpsertSimBatch", "UpsertPackage", "UpsertInventory", "UpsertBilling", "UpsertHealthReport",
		"DeleteCustomer", "DeleteSim", "DeletePackage",
	} {
		snap.On(m, mock.Anything).Return(errFor("snapshot")).Maybe()
	}
	snap.On("UpdateNodeStatus", mock.Anything, mock.Anything, mock.Anything).Return(errFor("snapshot")).Maybe()

	fact := &mocks.FactRepo{}
	for _, m := range []string{
		"AddPaymentEvent", "AddUsageEvent", "AddMetricSample", "AddAlarmEvent",
		"AddNodeStateEvent", "AddSiteStateEvent", "AddCustomerEvent", "AddSimEvent",
		"AddPackageEvent", "AddInventoryEvent",
	} {
		fact.On(m, mock.Anything).Return(errFor("fact")).Maybe()
	}
	fact.On("TransitionNodeState", mock.Anything, mock.Anything, mock.Anything).Return(errFor("fact")).Maybe()
	fact.On("TransitionSiteState", mock.Anything, mock.Anything, mock.Anything).Return(errFor("fact")).Maybe()
	fact.On("TransitionSimState", mock.Anything, mock.Anything, mock.Anything).Return(errFor("fact")).Maybe()
	fact.On("OpenCustomerPackageInterval", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errFor("fact")).Maybe()
	fact.On("CloseCustomerPackageInterval", mock.Anything, mock.Anything).Return(errFor("fact")).Maybe()

	st := &mocks.StateRepo{}
	st.On("MarkRollupDirty", mock.Anything).Return(errFor("state")).Maybe()
	st.On("UpsertRefreshState", mock.Anything).Return(nil).Maybe()
	st.On("SetRollupWatermark", mock.Anything, mock.Anything).Return(nil).Maybe()
	st.On("GetRefreshStates").Return(nil, nil).Maybe()
	st.On("GetRollupStates").Return(nil, nil).Maybe()

	return NewCollectorEventServer(testOrgName, ev, st, snap, fact)
}

func failureCases() []struct {
	event evt.EventId
	msg   proto.Message
} {
	return []struct {
		event evt.EventId
		msg   proto.Message
	}{
		{evt.EventPaymentSuccess, &epb.Payment{Id: "p1", Status: "success", AmountCents: 1, ItemType: "package"}},
		{evt.EventSubscriberCreate, &epb.EventSubscriberAdded{SubscriberId: "11111111-1111-1111-1111-111111111111"}},
		{evt.EventSubscriberDelete, &epb.EventSubscriberDeleted{SubscriberId: "11111111-1111-1111-1111-111111111111"}},
		{evt.EventSimAllocate, &epb.EventSimAllocation{Id: "sim-1", SubscriberId: "11111111-1111-1111-1111-111111111111"}},
		{evt.EventSimActivate, &epb.EventSimActivation{Id: "sim-1"}},
		{evt.EventSimAddPackage, &epb.EventSimAddPackage{Id: "sim-1"}},
		{evt.EventSimActivePackage, &epb.EventSimActivePackage{Id: "sim-1"}},
		{evt.EventSimRemovePackage, &epb.EventSimRemovePackage{Id: "sim-1"}},
		{evt.EventSimDelete, &epb.EventSimTermination{Id: "sim-1"}},
		{evt.EventSimsUpload, &epb.EventSimsUploaded{SimType: "ukama"}},
		{evt.EventPackageCreate, &epb.CreatePackageEvent{Uuid: "22222222-2222-2222-2222-222222222222"}},
		{evt.EventPackageUpdate, &epb.UpdatePackageEvent{Uuid: "22222222-2222-2222-2222-222222222222"}},
		{evt.EventPackageDelete, &epb.DeletePackageEvent{Uuid: "22222222-2222-2222-2222-222222222222"}},
		{evt.EventNetworkAdd, &epb.EventNetworkCreate{Id: "33333333-3333-3333-3333-333333333333"}},
		{evt.EventSiteCreate, &epb.EventAddSite{SiteId: "44444444-4444-4444-4444-444444444444", NetworkId: "33333333-3333-3333-3333-333333333333"}},
		{evt.EventSiteUpdate, &epb.EventUpdateSite{SiteId: "44444444-4444-4444-4444-444444444444"}},
		{evt.EventNodeCreate, &epb.EventRegistryNodeCreate{NodeId: "node-1"}},
		{evt.EventNodeUpdate, &epb.EventRegistryNodeUpdate{NodeId: "node-1"}},
		{evt.EventNodeAssign, &epb.EventRegistryNodeAssign{NodeId: "node-1", Site: "44444444-4444-4444-4444-444444444444"}},
		{evt.EventNodeRelease, &epb.EventRegistryNodeRelease{NodeId: "node-1"}},
		{evt.EventNodeOnline, &epb.NodeOnlineEvent{NodeId: "node-1"}},
		{evt.EventNodeOffline, &epb.NodeOfflineEvent{NodeId: "node-1"}},
		{evt.EventNodeStateTransition, &epb.NodeStateChangeEvent{NodeId: "node-1", State: "active"}},
		{evt.EventComponentsSync, &epb.EventInventoryNodeComponentAdd{Id: "comp-1"}},
		{evt.EventInvoiceGenerate, &epb.Report{Id: "inv-1"}},
	}
}

func TestDispatchHandlers_FailurePaths(t *testing.T) {
	for _, layer := range []string{"snapshot", "fact", "state"} {
		for _, tc := range failureCases() {
			t.Run(layer+"/"+string(evt.EventRoutingKey[tc.event]), func(t *testing.T) {
				s := failingServer(layer)

				am, _ := anypb.New(tc.msg)
				// Exercises the handler's failure branch; result is intentionally
				// ignored (recordFailure may return a nil response).
				_, _ = s.EventNotification(context.TODO(), &epb.Event{
					RoutingKey: msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[tc.event]),
					Msg:        am,
				})
			})
		}
	}
}
