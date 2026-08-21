/*
* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at https://mozilla.org/MPL/2.0/.
*
* Copyright (c) 2026-present, Ukama Inc.
 */

package reconciler

import (
	"fmt"

	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
)

const (
	SyncStatusInSync   = "in_sync"
	SyncStatusApplying = "applying"
	SyncStatusFault    = "fault"
)

func syncStatus(intent *db.SiteIntent, state *db.SiteState, flight *db.SiteIntentFlight) (string, string) {
	if intentMatchesState(intent, state) {
		return SyncStatusInSync, ReasonOK
	}

	if flight == nil {
		return SyncStatusApplying, "not applied yet"
	}

	switch flight.Status {
	case db.IntentFlightStatusFailed, db.IntentFlightStatusTimeout, db.IntentFlightStatusExpired:
		return SyncStatusFault, driftReason(intent, state, flight)
	default:
		return SyncStatusApplying, driftReason(intent, state, flight)
	}
}

func driftReason(intent *db.SiteIntent, state *db.SiteState, flight *db.SiteIntentFlight) string {
	retries := 0
	if flight != nil {
		retries = flight.RetryCount
	}

	if !serviceStateMatches(intent.DesiredService, stateService(state)) {
		return fmt.Sprintf("service desired=%s observed=%s attempts=%d",
			intent.DesiredService, stateService(state), retries)
	}

	return fmt.Sprintf("radio desired=%s observed=%s attempts=%d",
		intent.DesiredRadio, stateRadio(state), retries)
}

func (r *Reconciler) SyncStatus(siteID string) (string, string, error) {
	intent, err := r.getIntent(siteID)
	if err != nil {
		return "", "", err
	}

	state, err := r.states.Get(siteID)
	if err != nil {
		return "", "", err
	}

	flight, err := r.getFlight(intent)
	if err != nil {
		return "", "", err
	}

	status, reason := syncStatus(intent, state, flight)

	return status, reason, nil
}
