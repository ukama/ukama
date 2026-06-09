/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package refresh_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/analytics/collector/mocks"
	"github.com/ukama/ukama/systems/analytics/collector/pkg/clients"
	"github.com/ukama/ukama/systems/analytics/collector/pkg/refresh"
)

func newRefresher() (*refresh.Refresher, *mocks.RegistryClient, *mocks.SubscriberClient,
	*mocks.DataplanClient, *mocks.MetricsClient, *mocks.NodeClient, *mocks.InventoryClient, *mocks.BillingClient) {

	state := &mocks.StateRepo{}
	snap := &mocks.SnapshotRepo{}
	fact := &mocks.FactRepo{}

	state.On("UpsertRefreshState", mock.Anything).Return(nil).Maybe()
	state.On("MarkRollupDirty", mock.Anything).Return(nil).Maybe()
	state.On("SetRollupWatermark", mock.Anything, mock.Anything).Return(nil).Maybe()

	// snapshot + fact repos: everything succeeds.
	for _, m := range []string{
		"UpsertNetwork", "UpsertSite", "UpsertNode", "UpsertCustomer", "UpsertSim",
		"UpsertSimBatch", "UpsertPackage", "UpsertInventory", "UpsertBilling", "UpsertHealthReport",
	} {
		snap.On(m, mock.Anything).Return(nil).Maybe()
	}
	snap.On("UpdateNodeStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	snap.On("DeleteCustomer", mock.Anything).Return(nil).Maybe()
	snap.On("DeleteSim", mock.Anything).Return(nil).Maybe()
	snap.On("DeletePackage", mock.Anything).Return(nil).Maybe()

	for _, m := range []string{
		"AddPaymentEvent", "AddUsageEvent", "AddMetricSample", "AddAlarmEvent",
		"AddNodeStateEvent", "AddSiteStateEvent", "AddCustomerEvent", "AddSimEvent",
		"AddPackageEvent", "AddInventoryEvent",
	} {
		fact.On(m, mock.Anything).Return(nil).Maybe()
	}

	reg := &mocks.RegistryClient{}
	sub := &mocks.SubscriberClient{}
	dp := &mocks.DataplanClient{}
	met := &mocks.MetricsClient{}
	node := &mocks.NodeClient{}
	inv := &mocks.InventoryClient{}
	bill := &mocks.BillingClient{}

	id1 := "11111111-1111-1111-1111-111111111111"
	id2 := "22222222-2222-2222-2222-222222222222"

	reg.On("GetNetworks").Return([]clients.RegistryNetwork{{Id: id1, Name: "Net"}}, nil).Maybe()
	reg.On("GetSites").Return([]clients.RegistrySite{{Id: id2, NetworkId: id1, Latitude: "1.0", Longitude: "2.0"}}, nil).Maybe()
	sub.On("GetSubscribers").Return([]clients.SubscriberRecord{{SubscriberId: id1, NetworkId: id2}}, nil).Maybe()
	dp.On("GetPackages").Return([]clients.DataplanPackage{{Uuid: id1}}, nil).Maybe()
	met.On("GetLatestMetrics").Return([]clients.MetricValue{{}}, nil).Maybe()
	node.On("GetNodes").Return([]clients.NodeRecord{{Id: "node-1", SiteId: id2, NetworkId: id1}}, nil).Maybe()
	inv.On("GetComponents").Return([]clients.InventoryComponent{{}}, nil).Maybe()
	bill.On("GetAccount").Return(&clients.BillingAccount{}, nil).Maybe()

	r := refresh.NewRefresher(state, snap, fact, reg, sub, dp, met, node, inv, bill)
	return r, reg, sub, dp, met, node, inv, bill
}

func TestRefresher_AllSources(t *testing.T) {
	for _, source := range []string{
		"registry", "subscriber", "dataplan", "metrics", "node", "inventory", "billing",
	} {
		t.Run(source, func(t *testing.T) {
			r, _, _, _, _, _, _, _ := newRefresher()

			state, _ := r.Refresh(source)

			assert.NotNil(t, state)
			assert.Equal(t, source, state.Source)
		})
	}
}

func TestRefresher_UnknownSource(t *testing.T) {
	r, _, _, _, _, _, _, _ := newRefresher()

	state, err := r.Refresh("bogus")

	assert.Error(t, err)
	assert.NotNil(t, state)
	assert.Equal(t, "failed", state.Status)
}
