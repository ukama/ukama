/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package sql

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gorm.io/gorm"
)

const (
	opCreate = "create"
	opQuery  = "query"
	opUpdate = "update"
	opDelete = "delete"
	opRow    = "row"
	opRaw    = "raw"
)

const (
	statusOK    = "ok"
	statusError = "error"
)

// DbRequests counts database requests issued through GORM, by operation and outcome.
var DbRequests = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "ukama_db_requests_total",
	Help: "Total database requests issued through GORM, by operation and outcome.",
}, []string{"operation", "status"})

func init() {
	// Pre-create the series so an idle service reports 0 rather than nothing.
	for _, op := range []string{opCreate, opQuery, opUpdate, opDelete, opRow, opRaw} {
		for _, status := range []string{statusOK, statusError} {
			DbRequests.WithLabelValues(op, status)
		}
	}
}

// registerMetricsCallbacks registers the request counter on each GORM operation.
func registerMetricsCallbacks(g *gorm.DB) error {
	cb := g.Callback()

	registrations := []func() error{
		func() error {
			return cb.Create().After("gorm:create").Register("ukama:metrics_create", count(opCreate))
		},
		func() error {
			return cb.Query().After("gorm:query").Register("ukama:metrics_query", count(opQuery))
		},
		func() error {
			return cb.Update().After("gorm:update").Register("ukama:metrics_update", count(opUpdate))
		},
		func() error {
			return cb.Delete().After("gorm:delete").Register("ukama:metrics_delete", count(opDelete))
		},
		func() error {
			return cb.Row().After("gorm:row").Register("ukama:metrics_row", count(opRow))
		},
		func() error {
			return cb.Raw().After("gorm:raw").Register("ukama:metrics_raw", count(opRaw))
		},
	}

	for _, register := range registrations {
		if err := register(); err != nil {
			return err
		}
	}

	return nil
}

func count(op string) func(*gorm.DB) {
	return func(tx *gorm.DB) {
		status := statusOK

		// ErrRecordNotFound is a normal lookup outcome, not a failure.
		if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			status = statusError
		}

		DbRequests.WithLabelValues(op, status).Inc()
	}
}
