/*
* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at https://mozilla.org/MPL/2.0/.
*
* Copyright (c) 2026-present, Ukama Inc.
 */

package reconciler

import "github.com/ukama/ukama/systems/node/site-controller/pkg/db"

func expectedServiceState(desired string) string {
	if desired == StateOn {
		return StateRunning
	}
	return StateOff
}

func expectedRadioState(desired string) string {
	if desired == StateOn {
		return StateOn
	}
	return StateOff
}

func serviceStateMatches(desired, actual string) bool {
	return actual == expectedServiceState(desired)
}

func radioStateMatches(desired, actual string) bool {
	return actual == expectedRadioState(desired)
}

// intentMatchesState returns true when device-reported SiteState aligns with SiteIntent.
// Missing or unknown SiteState is treated as out of sync.
func intentMatchesState(intent *db.SiteIntent, state *db.SiteState) bool {
	if intent == nil {
		return false
	}
	return serviceStateMatches(intent.DesiredService, stateService(state)) &&
		radioStateMatches(intent.DesiredRadio, stateRadio(state))
}

func stateService(state *db.SiteState) string {
	if state == nil {
		return StateUnknown
	}
	return state.ServiceState
}

func stateRadio(state *db.SiteState) string {
	if state == nil {
		return StateUnknown
	}
	return state.RadioState
}
