/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Pure, framework-free aggregation logic for operation-lock status.
 *
 * Kept separate from the resolvers so it can be unit-tested without a
 * GraphQL context or network. All functions are stateless and take an
 * injectable `now` for deterministic tests.
 */
import { NODE_TYPE } from "../common/enums";
import {
  NodeOperationStatusDto,
  OperationDto,
  ResourceLockDto,
  SiteActionsDto,
} from "./resolvers/types";

/** Operation lifecycle states that mean the lock is (or is about to be) free. */
const TERMINAL_STATUSES = new Set([
  "success",
  "failed",
  "timeout",
  "cancelled",
]);

/** The operation system always keys locks as `node:<NODE_ID>`. */
export const nodeResourceKey = (nodeId: string): string => `node:${nodeId}`;

const isTerminal = (op?: OperationDto): boolean =>
  !!op && TERMINAL_STATUSES.has(op.status?.toLowerCase());

/**
 * A lease that has already elapsed is treated as free: the operation
 * system's sweeper (30s interval) will flip it to TIMEOUT shortly, and the
 * UI should not show a stuck spinner in the meantime.
 */
export const leaseExpired = (op: OperationDto | undefined, now: number): boolean => {
  if (!op?.leaseExpiresAt) return false;
  const t = Date.parse(op.leaseExpiresAt);
  return Number.isFinite(t) && t <= now;
};

/** A node is busy iff its lock is held by a non-terminal, unexpired operation. */
export const isLockBusy = (
  lock: ResourceLockDto | undefined,
  now: number = Date.now()
): boolean => {
  if (!lock || !lock.locked) return false;
  const op = lock.operation;
  if (isTerminal(op) || leaseExpired(op, now)) return false;
  return true;
};

/** The operation to surface for a node — only when genuinely busy. */
export const activeOperation = (
  lock: ResourceLockDto | undefined,
  now: number = Date.now()
): OperationDto | undefined => (isLockBusy(lock, now) ? lock?.operation : undefined);

/** A per-request read result for one node (lock resolved, or read failed). */
export interface NodeLockRead {
  id: string;
  type?: NODE_TYPE;
  lock?: ResourceLockDto;
  /** True when the lock read failed — fail-open (treated as not busy). */
  failed?: boolean;
}

export const toNodeStatus = (
  n: NodeLockRead,
  now: number = Date.now()
): NodeOperationStatusDto => ({
  nodeId: n.id,
  type: n.type,
  busy: n.failed ? false : isLockBusy(n.lock, now),
  operation: n.failed ? undefined : activeOperation(n.lock, now),
});

const OP_LABELS: Record<string, string> = {
  restartnode: "Restart",
  restartsite: "Site restart",
  togglerf: "RF change",
  togglerfstatus: "RF change",
  toggleservice: "Service change",
};

/** Short, human label for an operation type; falls back to the raw type. */
export const prettyOpType = (type?: string): string => {
  if (!type) return "An operation";
  return OP_LABELS[type.toLowerCase()] ?? type;
};

/** Reason string for a disabled action, derived from the blocking op. */
export const busyReason = (op?: OperationDto): string => {
  const who = op?.requestedBy ? ` by ${op.requestedBy}` : "";
  return `${prettyOpType(op?.type)} in progress${who}`;
};

/**
 * Per-action availability for the site page, derived from node roles:
 *  - restartSite depends on EVERY node (all must be idle)
 *  - rf         depends on the amplifier node (anode)
 *  - service    depends on the tower node (tnode)
 */
export const buildSiteActions = (
  statuses: NodeOperationStatusDto[]
): SiteActionsDto => {
  const tower = statuses.find((s) => s.type === NODE_TYPE.tnode);
  const amp = statuses.find((s) => s.type === NODE_TYPE.anode);
  const firstBusy = statuses.find((s) => s.busy);

  const restartSite: SiteActionsDto["restartSite"] = firstBusy
    ? { available: false, reason: busyReason(firstBusy.operation) }
    : { available: true };

  const roleAction = (
    node: NodeOperationStatusDto | undefined,
    missingReason: string
  ): SiteActionsDto["rf"] => {
    if (!node) return { available: false, reason: missingReason };
    if (node.busy) return { available: false, reason: busyReason(node.operation) };
    return { available: true };
  };

  return {
    restartSite,
    rf: roleAction(amp, "No amplifier node on this site"),
    service: roleAction(tower, "No tower node on this site"),
  };
};
