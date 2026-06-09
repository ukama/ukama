/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import {
  OperationDto,
  ResourceLockDto,
  StartOperationResponseDto,
} from "../resolvers/types";

/* The operation gateway emits snake_case JSON at the REST boundary. */
interface OperationRest {
  id: string;
  type: string;
  system: string;
  status: string;
  fencing_token: number;
  requested_by?: string;
  idempotency_key?: string;
  resource_key: string;
  lease_expires_at?: string;
  error?: string;
  started_at?: string;
  terminal_at?: string;
  created_at?: string;
}

const mapOperation = (op?: OperationRest | null): OperationDto | undefined => {
  if (!op) {
    return undefined;
  }
  return {
    id: op.id,
    type: op.type,
    system: op.system,
    status: op.status,
    fencingToken: op.fencing_token,
    requestedBy: op.requested_by,
    idempotencyKey: op.idempotency_key,
    resourceKey: op.resource_key,
    leaseExpiresAt: op.lease_expires_at,
    error: op.error,
    startedAt: op.started_at,
    terminalAt: op.terminal_at,
    createdAt: op.created_at,
  };
};

export const mapStartOperation = (res: {
  operation?: OperationRest;
  conflicting_operation?: OperationRest;
}): StartOperationResponseDto => ({
  operation: mapOperation(res.operation),
  conflictingOperation: mapOperation(res.conflicting_operation),
});

export const mapGetOperation = (res: {
  operation?: OperationRest;
}): OperationDto | undefined => mapOperation(res.operation);

export const mapResourceLock = (res: {
  locked: boolean;
  operation?: OperationRest;
}): ResourceLockDto => ({
  locked: res.locked,
  operation: mapOperation(res.operation),
});
