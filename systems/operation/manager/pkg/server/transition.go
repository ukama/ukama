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
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
)

type Repos struct {
	Operations db.OperationRepo
	Locks      db.LockRepo
	Audits     db.AuditRepo
}

type refusal struct {
	op     db.Operation
	prefix string
	reason string
	cause  error
}

func (r *refusal) Error() string {
	if r.cause != nil {
		return r.cause.Error()
	}
	return r.prefix + "transition"
}

func reject(op *db.Operation, cause error) *refusal {
	return &refusal{op: *op, prefix: "rejected_", cause: cause}
}

func ignore(op *db.Operation, cause error) *refusal {
	return &refusal{op: *op, prefix: "ignored_", cause: cause}
}

func (r Repos) markRunning(ctx context.Context, id uuid.UUID, fencingToken uint64,
	actor string) (*db.Operation, error) {
	op, err := r.Operations.Update(ctx, id, func(op *db.Operation, tx *gorm.DB) error {
		now := time.Now().UTC()
		var ref *refusal
		switch {
		case op.FencingToken != fencingToken:
			ref = reject(op, fmt.Errorf("fencing token mismatch (op=%d, given=%d): %w",
				op.FencingToken, fencingToken, store.ErrPrecondition))
		case op.Status.IsTerminal():
			ref = reject(op, fmt.Errorf("operation %s is already %s: %w",
				id, op.Status, store.ErrPrecondition))
		case !op.LeaseExpiresAt.After(now):
			ref = reject(op, fmt.Errorf("operation %s lease expired at %s: %w",
				id, op.LeaseExpiresAt.Format(time.RFC3339), store.ErrPrecondition))
		case op.Status == db.OperationRunning:
			return ignore(op, nil)
		}
		if ref != nil {
			ref.reason = ref.cause.Error()
			return ref
		}

		op.Status = db.OperationRunning
		op.StartedAt = &now
		return r.Audits.AddInTx(tx, &db.OperationAudit{
			OperationId: op.Id, ResourceKey: op.ResourceKey, Event: "running",
			Actor: actor, At: now,
		})
	})
	return r.settle(ctx, op, err, db.OperationAudit{Event: "running", Actor: actor})
}

func (r Repos) terminate(ctx context.Context, id uuid.UUID, fencingToken uint64,
	status db.OperationStatus, audit db.OperationAudit,
	opErr string) (*db.Operation, error) {
	if !status.IsTerminal() {
		return nil, fmt.Errorf("terminate: %s is not a terminal status", status)
	}

	op, err := r.Operations.Update(ctx, id, func(op *db.Operation, tx *gorm.DB) error {
		if op.FencingToken != fencingToken {
			return reject(op, fmt.Errorf("fencing token mismatch (op=%d, given=%d): %w",
				op.FencingToken, fencingToken, store.ErrPrecondition))
		}
		if op.Status.IsTerminal() {
			var cause error
			if op.Status != status {
				cause = fmt.Errorf("operation %s is already %s, cannot move to %s: %w",
					id, op.Status, status, store.ErrPrecondition)
			}
			return ignore(op, cause)
		}
		if status == db.OperationSuccess && op.Status == db.OperationPending {
			return reject(op, fmt.Errorf("operation %s is pending, cannot move to %s: %w",
				id, status, store.ErrPrecondition))
		}

		now := time.Now().UTC()
		op.Status = status
		op.TerminalAt = &now
		if opErr != "" {
			op.Error = opErr
		}
		if err := r.Locks.Release(tx, op); err != nil {
			return err
		}
		final := audit
		final.OperationId = op.Id
		final.ResourceKey = op.ResourceKey
		final.At = now
		return r.Audits.AddInTx(tx, &final)
	})
	return r.settle(ctx, op, err, audit)
}

func (r Repos) settle(ctx context.Context, op *db.Operation, err error,
	audit db.OperationAudit) (*db.Operation, error) {
	var ref *refusal
	if !errors.As(err, &ref) {
		return op, err
	}

	audit.OperationId = ref.op.Id
	audit.ResourceKey = ref.op.ResourceKey
	audit.Event = ref.prefix + audit.Event
	if ref.reason != "" {
		audit.Reason = ref.reason
	}
	if aerr := r.Audits.Add(ctx, &audit); aerr != nil {
		return nil, aerr
	}
	if ref.cause != nil {
		return nil, ref.cause
	}
	return &ref.op, nil
}
