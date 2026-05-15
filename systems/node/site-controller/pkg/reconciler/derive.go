/*
* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at https://mozilla.org/MPL/2.0/.
*
* Copyright (c) 2026-present, Ukama Inc.
 */

package reconciler

import "github.com/ukama/ukama/systems/node/site-controller/pkg/db"

func derive(intent *db.SiteIntent) *db.SiteState {
	state := &db.SiteState{SiteID: intent.SiteID, PowerState: StateHealthy, ServiceState: StateUnknown, RadioState: StateUnknown, AccessState: StateUnavailable, Reason: ReasonOK}
	if intent.DesiredService == StateOn {
		state.ServiceState = StateRunning
	} else {
		state.ServiceState = StateOff
	}
	if intent.DesiredRadio == StateOn {
		state.RadioState = StateOn
	} else {
		state.RadioState = StateOff
	}
	if intent.DesiredService == StateOn && intent.DesiredRadio == StateOn {
		state.AccessState = StateAvailable
	} else if intent.DesiredRadio == StateOff {
		state.Reason = "radio_off"
	} else if intent.DesiredService == StateOff {
		state.Reason = "service_off"
	}
	return state
}
