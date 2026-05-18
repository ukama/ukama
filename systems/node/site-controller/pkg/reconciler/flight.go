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

	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
)

func (r *Reconciler) getFlight(intent *db.SiteIntent) (*db.SiteIntentFlight, error) {
	if r.flights == nil || intent == nil || intent.ID == uuid.Nil {
		return nil, nil
	}
	return r.flights.GetBySiteIntentID(intent.ID)
}

func (r *Reconciler) resetIntentReconcile(intent *db.SiteIntent) error {
	return r.flights.Upsert(&db.SiteIntentFlight{
		SiteIntentID: intent.ID,
		Status:       db.IntentFlightStatusPending,
		RetryCount:   0,
		ExpiresAt:    time.Now().UTC().Add(r.flightTTL),
	})
}

func (r *Reconciler) markFlightStatus(intent *db.SiteIntent, status string) error {
	if r.flights == nil || intent == nil || intent.ID == uuid.Nil {
		return nil
	}
	flight, err := r.getFlight(intent)
	if err != nil {
		return err
	}
	expiresAt := time.Now().UTC().Add(r.flightTTL)
	retryCount := 0
	if flight != nil {
		expiresAt = flight.ExpiresAt
		retryCount = flight.RetryCount
	}
	return r.flights.Upsert(&db.SiteIntentFlight{
		SiteIntentID: intent.ID,
		Status:       status,
		RetryCount:   retryCount,
		ExpiresAt:    expiresAt,
	})
}

func (r *Reconciler) saveFlight(flight *db.SiteIntentFlight) error {
	if flight == nil {
		return nil
	}
	return r.flights.Upsert(flight)
}

func (r *Reconciler) flightExpired(flight *db.SiteIntentFlight) bool {
	return flight != nil && !flight.ExpiresAt.IsZero() && time.Now().UTC().After(flight.ExpiresAt)
}

func (r *Reconciler) flightRetriesExhausted(flight *db.SiteIntentFlight) bool {
	return flight != nil && r.maxRetries > 0 && flight.RetryCount >= r.maxRetries
}

func (r *Reconciler) flightDue(flight *db.SiteIntentFlight) bool {
	if flight == nil || flight.Status != db.IntentFlightStatusPending {
		return true
	}
	return time.Since(flight.UpdatedAt) >= r.reconcileInterval
}
