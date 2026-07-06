/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { NODE_TYPE } from "../common/enums";
import {
  buildSiteActions,
  isLockBusy,
  leaseExpired,
  nodeResourceKey,
  toNodeStatus,
} from "./logic";
import { NodeOperationStatusDto, OperationDto, ResourceLockDto } from "./resolvers/types";

const NOW = Date.parse("2026-07-06T12:00:00Z");

const op = (over: Partial<OperationDto> = {}): OperationDto => ({
  id: "op-1",
  type: "RestartNode",
  system: "node",
  status: "running",
  fencingToken: 1,
  resourceKey: "node:n1",
  requestedBy: "alan",
  leaseExpiresAt: new Date(NOW + 60_000).toISOString(),
  ...over,
});

const lock = (over: Partial<ResourceLockDto> = {}): ResourceLockDto => ({
  locked: true,
  operation: op(),
  ...over,
});

const status = (over: Partial<NodeOperationStatusDto>): NodeOperationStatusDto => ({
  nodeId: "n",
  busy: false,
  ...over,
});

describe("nodeResourceKey", () => {
  it("prefixes with node:", () => {
    expect(nodeResourceKey("uk-123")).toBe("node:uk-123");
  });
});

describe("isLockBusy", () => {
  it("false when not locked or missing", () => {
    expect(isLockBusy(undefined, NOW)).toBe(false);
    expect(isLockBusy({ locked: false }, NOW)).toBe(false);
  });
  it("true when locked with a live non-terminal op", () => {
    expect(isLockBusy(lock(), NOW)).toBe(true);
  });
  it("false when op is terminal (any case)", () => {
    expect(isLockBusy(lock({ operation: op({ status: "SUCCESS" }) }), NOW)).toBe(false);
    expect(isLockBusy(lock({ operation: op({ status: "failed" }) }), NOW)).toBe(false);
  });
  it("false when lease already expired (sweeper will reclaim)", () => {
    const expired = op({ leaseExpiresAt: new Date(NOW - 1000).toISOString() });
    expect(isLockBusy(lock({ operation: expired }), NOW)).toBe(false);
  });
});

describe("leaseExpired", () => {
  it("false when no lease provided", () => {
    expect(leaseExpired(op({ leaseExpiresAt: undefined }), NOW)).toBe(false);
  });
  it("true only once the lease timestamp has passed", () => {
    expect(leaseExpired(op({ leaseExpiresAt: new Date(NOW + 1).toISOString() }), NOW)).toBe(false);
    expect(leaseExpired(op({ leaseExpiresAt: new Date(NOW - 1).toISOString() }), NOW)).toBe(true);
  });
});

describe("toNodeStatus", () => {
  it("failed read fails open (idle, no op)", () => {
    const s = toNodeStatus({ id: "n1", type: NODE_TYPE.tnode, failed: true }, NOW);
    expect(s).toEqual({ nodeId: "n1", type: NODE_TYPE.tnode, busy: false, operation: undefined });
  });
  it("busy node surfaces its operation", () => {
    const s = toNodeStatus({ id: "n1", type: NODE_TYPE.anode, lock: lock() }, NOW);
    expect(s.busy).toBe(true);
    expect(s.operation?.id).toBe("op-1");
  });
});

describe("buildSiteActions — independent async release", () => {
  it("all idle → everything available", () => {
    const a = buildSiteActions([
      status({ nodeId: "t", type: NODE_TYPE.tnode }),
      status({ nodeId: "a", type: NODE_TYPE.anode }),
    ]);
    expect(a.restartSite.available).toBe(true);
    expect(a.rf.available).toBe(true);
    expect(a.service.available).toBe(true);
  });

  it("amplifier busy → RF and restartSite locked, service stays available", () => {
    const a = buildSiteActions([
      status({ nodeId: "t", type: NODE_TYPE.tnode }),
      status({ nodeId: "a", type: NODE_TYPE.anode, busy: true, operation: op({ type: "ToggleRF", requestedBy: "sam" }) }),
    ]);
    expect(a.rf.available).toBe(false);
    expect(a.rf.reason).toContain("sam");
    expect(a.restartSite.available).toBe(false);
    expect(a.service.available).toBe(true);
  });

  it("tower busy → service and restartSite locked, RF stays available", () => {
    const a = buildSiteActions([
      status({ nodeId: "t", type: NODE_TYPE.tnode, busy: true, operation: op({ type: "ToggleService" }) }),
      status({ nodeId: "a", type: NODE_TYPE.anode }),
    ]);
    expect(a.service.available).toBe(false);
    expect(a.rf.available).toBe(true);
    expect(a.restartSite.available).toBe(false);
  });

  it("missing role node → action unavailable with a clear reason", () => {
    const a = buildSiteActions([status({ nodeId: "t", type: NODE_TYPE.tnode })]);
    expect(a.rf.available).toBe(false);
    expect(a.rf.reason).toBe("No amplifier node on this site");
  });
});
