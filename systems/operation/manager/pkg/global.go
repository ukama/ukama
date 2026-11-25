/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package pkg

const (
	SystemName  = "operation"
	ServiceName = "manager"
)

const (
	AuditLockAcquired = "lock_acquired"
	AuditRunning      = "running"
	AuditCompleted    = "completed"
	AuditFailed       = "failed"
	AuditTimeout      = "timeout"
	AuditForceUnlock  = "force_unlock"
	AuditDeleted      = "deleted"
	AuditRejected     = "rejected_"
	AuditIgnored      = "ignored_"

	ReasonLeaseExpired = "lease expired"
)

var (
	IsDebugMode bool   = false
	InstanceId  string = ServiceName + "-debug"
)
