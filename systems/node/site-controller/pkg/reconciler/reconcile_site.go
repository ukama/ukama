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
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
)

// ReconcileSite drives service/radio actions until SiteState (device-reported, external)
// matches SiteIntent. This service never writes SiteState.
func (r *Reconciler) ReconcileSite(ctx context.Context, siteID string, force bool) error {
	intent, err := r.getIntent(siteID)
	if err != nil {
		return err
	}
	if intent.ID == uuid.Nil {
		return nil
	}

	flight, err := r.getFlight(intent)
	if err != nil {
		return err
	}
	if flight != nil && flight.IsTerminal() {
		return nil
	}

	state, err := r.states.Get(siteID)
	if err != nil {
		return err
	}
	if intentMatchesState(intent, state) {
		return r.markFlightStatus(intent, db.IntentFlightStatusSucceeded)
	}

	if r.flightExpired(flight) {
		log.Warnf("site-controller: flight for intent %s site %s expired", intent.ID, siteID)
		return r.markFlightStatus(intent, db.IntentFlightStatusExpired)
	}

	if r.flightRetriesExhausted(flight) {
		retries := 0
		if flight != nil {
			retries = flight.RetryCount
		}
		log.Warnf("site-controller: flight for intent %s site %s timed out after %d retries", intent.ID, siteID, retries)
		return r.markFlightStatus(intent, db.IntentFlightStatusTimeout)
	}

	if !force && !r.flightDue(flight) {
		return nil
	}

	if err := r.applyIntentState(ctx, intent, state); err != nil {
		return r.recordFailedAttempt(intent, flight, err)
	}

	return r.finishReconcileAttempt(siteID, intent, flight)
}

func (r *Reconciler) finishReconcileAttempt(siteID string, intent *db.SiteIntent, flight *db.SiteIntentFlight) error {
	state, err := r.states.Get(siteID)
	if err != nil {
		return err
	}
	if intentMatchesState(intent, state) {
		return r.markFlightStatus(intent, db.IntentFlightStatusSucceeded)
	}

	if flight == nil {
		flight = &db.SiteIntentFlight{SiteIntentID: intent.ID}
	}
	flight.RetryCount++
	if err := r.saveFlight(flight); err != nil {
		return err
	}

	if r.flightExpired(flight) {
		return r.markFlightStatus(intent, db.IntentFlightStatusExpired)
	}
	if r.flightRetriesExhausted(flight) {
		return r.markFlightStatus(intent, db.IntentFlightStatusTimeout)
	}
	return fmt.Errorf("site %s intent/state still out of sync (retry %d/%d)", siteID, flight.RetryCount, r.maxRetries)
}

// ResyncSite re-arms reconciliation for a site whose device-reported state has drifted
// from its intent. Once a site reaches its desired state its flight is parked in a
// terminal status, so the periodic worker alone would never correct drift introduced
// later, such as a node restart bringing service back up disabled.
func (r *Reconciler) ResyncSite(ctx context.Context, siteID string) error {
	if siteID == "" {
		return fmt.Errorf("site id is required")
	}

	intent, err := r.getIntent(siteID)
	if err != nil {
		return err
	}
	if intent.ID == uuid.Nil {
		return nil
	}

	state, err := r.states.Get(siteID)
	if err != nil {
		return err
	}
	if intentMatchesState(intent, state) {
		return nil
	}

	log.Infof("site-controller: site %s drifted from intent (service desired=%s observed=%s), re-arming reconcile",
		siteID, intent.DesiredService, stateService(state))

	if err := r.resetIntentReconcile(intent); err != nil {
		return fmt.Errorf("reset reconcile: %w", err)
	}

	return r.ReconcileSite(ctx, siteID, true)
}

func (r *Reconciler) applyIntentState(ctx context.Context, intent *db.SiteIntent, state *db.SiteState) error {
	siteID := intent.SiteID
	radioMismatch := !radioStateMatches(intent.DesiredRadio, stateRadio(state))
	serviceMismatch := !serviceStateMatches(intent.DesiredService, stateService(state))
	if !radioMismatch && !serviceMismatch {
		return nil
	}

	if intent.DesiredService == StateOn && serviceMismatch {
		if err := r.ensureCriticalPoe(ctx, siteID); err != nil {
			return fmt.Errorf("ensure critical poe: %w", err)
		}
	}

	// Radio and service are independent actuators, so both are attempted even when one
	// fails. Returning early on the radio error would leave a site stuck without service
	// whenever its radio cannot be resolved.
	var errs []error

	if radioMismatch {
		if err := r.applyRadio(ctx, siteID, intent.DesiredRadio); err != nil {
			errs = append(errs, fmt.Errorf("radio action: %w", err))
		}
	}
	if serviceMismatch {
		if err := r.applyService(ctx, siteID, intent.DesiredService); err != nil {
			errs = append(errs, fmt.Errorf("service action: %w", err))
		}
	}

	return errors.Join(errs...)
}
