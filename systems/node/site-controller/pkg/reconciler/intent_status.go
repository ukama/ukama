/*
* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at https://mozilla.org/MPL/2.0/.
*
* Copyright (c) 2026-present, Ukama Inc.
 */

package reconciler

import (
	"time"

	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
)

func (r *Reconciler) recordFailedAttempt(intent *db.SiteIntent, flight *db.SiteIntentFlight, applyErr error) error {
	if flight == nil {
		flight = &db.SiteIntentFlight{
			SiteIntentID: intent.ID,
			Status:       db.IntentFlightStatusPending,
			ExpiresAt:    time.Now().UTC().Add(r.flightTTL),
		}
	}
	flight.RetryCount++
	if applyErr != nil && r.flightRetriesExhausted(flight) {
		flight.Status = db.IntentFlightStatusFailed
		if err := r.saveFlight(flight); err != nil {
			return err
		}
		return applyErr
	}
	flight.Status = db.IntentFlightStatusPending
	if err := r.saveFlight(flight); err != nil {
		return err
	}
	if applyErr != nil {
		return applyErr
	}
	return nil
}
