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

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

func newTestEventServer() *CollectorEventServer {
	return NewCollectorEventServer(testOrgName, newStubEventRepo(), newStubStateRepo(),
		newStubSnapshotRepo(), &stubFactRepo{})
}

// TestDispatchAllRoutingKeys fires one event per supported routing key so the
// dispatch switch and every handler body executes.
func TestDispatchAllRoutingKeys(t *testing.T) {
	cases := []struct {
		event evt.EventId
		msg   proto.Message
	}{
		{evt.EventPaymentSuccess, &epb.Payment{Id: "pay-1", Status: "success", AmountCents: 100, ItemType: "package"}},
		{evt.EventPaymentFailed, &epb.Payment{Id: "pay-2", Status: "failed", AmountCents: 100, ItemType: "package"}},
		{evt.EventSubscriberCreate, &epb.EventSubscriberAdded{SubscriberId: "11111111-1111-1111-1111-111111111111", Name: "Jane"}},
		{evt.EventSubscriberUpdate, &epb.EventSubscriberAdded{SubscriberId: "11111111-1111-1111-1111-111111111111", Name: "Jane"}},
		{evt.EventSubscriberDelete, &epb.EventSubscriberDeleted{SubscriberId: "11111111-1111-1111-1111-111111111111"}},
		{evt.EventSimAllocate, &epb.EventSimAllocation{Id: "sim-1", SubscriberId: "11111111-1111-1111-1111-111111111111", Iccid: "8910"}},
		{evt.EventSimActivate, &epb.EventSimActivation{Id: "sim-1", Iccid: "8910"}},
		{evt.EventSimAddPackage, &epb.EventSimAddPackage{Id: "sim-1"}},
		{evt.EventSimActivePackage, &epb.EventSimActivePackage{Id: "sim-1"}},
		{evt.EventSimRemovePackage, &epb.EventSimRemovePackage{Id: "sim-1"}},
		{evt.EventSimDelete, &epb.EventSimTermination{Id: "sim-1"}},
		{evt.EventSimsUpload, &epb.EventSimsUploaded{SimType: "ukama"}},
		{evt.EventPackageCreate, &epb.CreatePackageEvent{Uuid: "22222222-2222-2222-2222-222222222222"}},
		{evt.EventPackageUpdate, &epb.UpdatePackageEvent{Uuid: "22222222-2222-2222-2222-222222222222"}},
		{evt.EventPackageDelete, &epb.DeletePackageEvent{Uuid: "22222222-2222-2222-2222-222222222222"}},
		{evt.EventNetworkAdd, &epb.EventNetworkCreate{Id: "33333333-3333-3333-3333-333333333333", Name: "Net"}},
		{evt.EventSiteCreate, &epb.EventAddSite{SiteId: "44444444-4444-4444-4444-444444444444", NetworkId: "33333333-3333-3333-3333-333333333333"}},
		{evt.EventSiteUpdate, &epb.EventUpdateSite{SiteId: "44444444-4444-4444-4444-444444444444"}},
		{evt.EventNodeCreate, &epb.EventRegistryNodeCreate{NodeId: "node-1", Name: "Node"}},
		{evt.EventNodeUpdate, &epb.EventRegistryNodeUpdate{NodeId: "node-1"}},
		{evt.EventNodeAssign, &epb.EventRegistryNodeAssign{NodeId: "node-1", Site: "44444444-4444-4444-4444-444444444444", Network: "33333333-3333-3333-3333-333333333333"}},
		{evt.EventNodeRelease, &epb.EventRegistryNodeRelease{NodeId: "node-1"}},
		{evt.EventNodeOnline, &epb.NodeOnlineEvent{NodeId: "node-1"}},
		{evt.EventNodeOffline, &epb.NodeOfflineEvent{NodeId: "node-1"}},
		{evt.EventNodeStateTransition, &epb.NodeStateChangeEvent{NodeId: "node-1", State: "active"}},
		{evt.EventHealthReportStore, &epb.HealthReportEvent{Id: "hr-1"}},
		{evt.EventComponentsSync, &epb.EventInventoryNodeComponentAdd{Id: "comp-1"}},
		{evt.EventInvoiceGenerate, &epb.Report{Id: "inv-1"}},
	}

	for _, tc := range cases {
		t.Run(string(evt.EventRoutingKey[tc.event]), func(t *testing.T) {
			s := newTestEventServer()

			anyMsg, err := anypb.New(tc.msg)
			assert.NoError(t, err)

			resp, err := s.EventNotification(context.TODO(), &epb.Event{
				RoutingKey: msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[tc.event]),
				Msg:        anyMsg,
			})

			assert.NoError(t, err)
			assert.NotNil(t, resp)
		})
	}
}

func TestDispatch_MalformedMessage(t *testing.T) {
	s := newTestEventServer()

	// wrong payload type for a payment routing key triggers recordMalformed
	anyMsg, _ := anypb.New(&epb.EventNetworkCreate{Id: "x"})

	resp, err := s.EventNotification(context.TODO(), &epb.Event{
		RoutingKey: msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventPaymentSuccess]),
		Msg:        anyMsg,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
