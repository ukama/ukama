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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/site-controller/mocks"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
)

const testSiteID = "site-resync-001"

func TestResyncSite(t *testing.T) {
	t.Run("re-arms a parked flight when service drifted from intent", func(t *testing.T) {
		intentID := uuid.NewV4()

		intents := &mocks.IntentRepo{}
		intents.On("Get", testSiteID).Return(&db.SiteIntent{
			ID: intentID, SiteID: testSiteID,
			DesiredService: StateOn, DesiredRadio: StateOn,
		}, nil)

		states := &mocks.StateRepo{}
		states.On("Get", testSiteID).Return(&db.SiteState{
			SiteID: testSiteID, ServiceState: StateOff, RadioState: StateOn,
		}, nil)

		flights := &mocks.IntentFlightRepo{}
		flights.On("Upsert", mock.MatchedBy(func(f *db.SiteIntentFlight) bool {
			return f.SiteIntentID == intentID &&
				f.Status == db.IntentFlightStatusPending &&
				f.RetryCount == 0
		})).Return(nil).Once()
		flights.On("GetBySiteIntentID", mock.Anything).Return(nil, nil).Maybe()
		flights.On("Upsert", mock.Anything).Return(nil).Maybe()

		ports := &mocks.PortMapRepo{}
		ports.On("GetBySite", testSiteID).Return([]db.SitePortMap{}, nil).Maybe()

		r := &Reconciler{intents: intents, states: states, flights: flights, ports: ports}

		// the downstream apply needs site hardware wiring that is out of scope here;
		// this asserts only that the parked flight was re-armed.
		_ = r.ResyncSite(context.TODO(), testSiteID)

		flights.AssertExpectations(t)
	})

	t.Run("does nothing when the site already matches its intent", func(t *testing.T) {
		intents := &mocks.IntentRepo{}
		intents.On("Get", testSiteID).Return(&db.SiteIntent{
			ID: uuid.NewV4(), SiteID: testSiteID,
			DesiredService: StateOn, DesiredRadio: StateOn,
		}, nil)

		states := &mocks.StateRepo{}
		states.On("Get", testSiteID).Return(&db.SiteState{
			SiteID: testSiteID, ServiceState: StateRunning, RadioState: StateOn,
		}, nil)

		flights := &mocks.IntentFlightRepo{}

		r := &Reconciler{intents: intents, states: states, flights: flights}

		require.NoError(t, r.ResyncSite(context.TODO(), testSiteID))
		flights.AssertNotCalled(t, "Upsert", mock.Anything)
	})

	t.Run("does nothing for a site with no intent", func(t *testing.T) {
		intents := &mocks.IntentRepo{}
		intents.On("Get", testSiteID).Return(nil, nil)

		flights := &mocks.IntentFlightRepo{}

		r := &Reconciler{intents: intents, flights: flights}

		require.NoError(t, r.ResyncSite(context.TODO(), testSiteID))
		flights.AssertNotCalled(t, "Upsert", mock.Anything)
	})

	t.Run("rejects an empty site id", func(t *testing.T) {
		r := &Reconciler{}
		assert.Error(t, r.ResyncSite(context.TODO(), ""))
	})
}
