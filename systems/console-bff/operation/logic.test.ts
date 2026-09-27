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
import {
  NodeOperationStatusDto,
  OperationDto,
  ResourceLockDto,
} from "./resolvers/types";

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

const status = (
  over: Partial<NodeOperationStatusDto>
): NodeOperationStatusDto => ({
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
    expect(
      isLockBusy(lock({ operation: op({ status: "SUCCESS" }) }), NOW)
    ).toBe(false);
    expect(isLockBusy(lock({ operation: op({ status: "failed" }) }), NOW)).toBe(
      false
    );
  });
  it("true when lease expired but the manager has not released the lock", () => {
    const expired = op({ leaseExpiresAt: new Date(NOW - 1000).toISOString() });
    expect(isLockBusy(lock({ operation: expired }), NOW)).toBe(true);
  });
});

describe("leaseExpired", () => {
  it("false when no lease provided", () => {
    expect(leaseExpired(op({ leaseExpiresAt: undefined }), NOW)).toBe(false);
  });
  it("true only once the lease timestamp has passed", () => {
    expect(
      leaseExpired(op({ leaseExpiresAt: new Date(NOW + 1).toISOString() }), NOW)
    ).toBe(false);
    expect(
      leaseExpired(op({ leaseExpiresAt: new Date(NOW - 1).toISOString() }), NOW)
    ).toBe(true);
  });
});

describe("toNodeStatus", () => {
  it("failed read fails open (idle, no op)", () => {
    const s = toNodeStatus(
      { id: "n1", type: NODE_TYPE.tnode, failed: true },
      NOW
    );
    expect(s).toEqual({
      nodeId: "n1",
      type: NODE_TYPE.tnode,
      busy: false,
      operation: undefined,
    });
  });
  it("busy node surfaces its operation", () => {
    const s = toNodeStatus(
      { id: "n1", type: NODE_TYPE.anode, lock: lock() },
      NOW
    );
    expect(s.busy).toBe(true);
    expect(s.operation?.id).toBe("op-1");
  });
});

describe("buildSiteActions — site-wide busy state", () => {
  it("keeps all actions unavailable after lease expiry until the manager releases it", () => {
    const expired = op({ leaseExpiresAt: new Date(NOW - 1000).toISOString() });
    const tower = toNodeStatus(
      { id: "t", type: NODE_TYPE.tnode, lock: lock({ operation: expired }) },
      NOW
    );
    const amp = status({ nodeId: "a", type: NODE_TYPE.anode });
    const waiting = buildSiteActions([tower, amp]);

    expect(tower.operation?.id).toBe(expired.id);
    expect(waiting.service.available).toBe(false);
    expect(waiting.restartSite.available).toBe(false);
    expect(waiting.rf.available).toBe(false);

    const released = toNodeStatus(
      { id: "t", type: NODE_TYPE.tnode, lock: { locked: false } },
      NOW
    );
    const ready = buildSiteActions([released, amp]);
    expect(released.operation).toBe(undefined);
    expect(ready.service.available).toBe(true);
    expect(ready.restartSite.available).toBe(true);
    expect(ready.rf.available).toBe(true);
  });

  it("all idle → everything available", () => {
    const a = buildSiteActions([
      status({ nodeId: "t", type: NODE_TYPE.tnode }),
      status({ nodeId: "a", type: NODE_TYPE.anode }),
      status({ nodeId: "c", type: NODE_TYPE.cnode }),
    ]);
    expect(a.restartSite.available).toBe(true);
    expect(a.rf.available).toBe(true);
    expect(a.service.available).toBe(true);
  });

  it("amplifier busy → all site actions locked", () => {
    const a = buildSiteActions([
      status({ nodeId: "t", type: NODE_TYPE.tnode }),
      status({
        nodeId: "a",
        type: NODE_TYPE.anode,
        busy: true,
        operation: op({ type: "ToggleRF", requestedBy: "sam" }),
      }),
    ]);
    expect(a.rf.available).toBe(false);
    expect(a.rf.reason).toContain("sam");
    expect(a.restartSite.available).toBe(false);
    expect(a.service.available).toBe(false);
    expect(a.service.reason).toBe(a.restartSite.reason);
  });

  it("tower busy → all site actions locked", () => {
    const a = buildSiteActions([
      status({
        nodeId: "t",
        type: NODE_TYPE.tnode,
        busy: true,
        operation: op({ type: "ToggleService" }),
      }),
      status({ nodeId: "a", type: NODE_TYPE.anode }),
    ]);
    expect(a.service.available).toBe(false);
    expect(a.rf.available).toBe(false);
    expect(a.restartSite.available).toBe(false);
    expect(a.rf.reason).toBe(a.restartSite.reason);
  });

  it("controller updating → all site actions locked with the same reason", () => {
    const a = buildSiteActions([
      status({ nodeId: "t", type: NODE_TYPE.tnode }),
      status({ nodeId: "a", type: NODE_TYPE.anode }),
      status({
        nodeId: "c",
        type: NODE_TYPE.cnode,
        busy: true,
        operation: op({ type: "UpdateSoftware", requestedBy: "software" }),
      }),
    ]);
    const blocked = {
      available: false,
      reason: "UpdateSoftware in progress by software",
    };
    expect(a.restartSite).toEqual(blocked);
    expect(a.rf).toEqual(blocked);
    expect(a.service).toEqual(blocked);
  });

  it("keeps all actions locked until every busy node is released", () => {
    const tower = status({ nodeId: "t", type: NODE_TYPE.tnode, busy: true });
    const amp = status({ nodeId: "a", type: NODE_TYPE.anode });
    const controller = status({ nodeId: "c", type: NODE_TYPE.cnode, busy: true });
    const busy = buildSiteActions([tower, amp, controller]);
    expect(busy.restartSite.available).toBe(false);
    expect(busy.rf.available).toBe(false);
    expect(busy.service.available).toBe(false);

    const waiting = buildSiteActions([{ ...tower, busy: false }, amp, controller]);
    expect(waiting.restartSite.available).toBe(false);
    expect(waiting.rf.available).toBe(false);
    expect(waiting.service.available).toBe(false);

    const ready = buildSiteActions([
      { ...tower, busy: false },
      amp,
      { ...controller, busy: false },
    ]);
    expect(ready.restartSite.available).toBe(true);
    expect(ready.rf.available).toBe(true);
    expect(ready.service.available).toBe(true);
  });

  it("missing role node → action unavailable with a clear reason", () => {
    const a = buildSiteActions([
      status({ nodeId: "t", type: NODE_TYPE.tnode }),
    ]);
    expect(a.rf.available).toBe(false);
    expect(a.rf.reason).toBe("No amplifier node on this site");
    const b = buildSiteActions([
      status({ nodeId: "a", type: NODE_TYPE.anode }),
    ]);
    expect(b.service.available).toBe(false);
    expect(b.service.reason).toBe("No tower node on this site");
  });

  it("preserves the missing-role reason while another node is busy", () => {
    const a = buildSiteActions([
      status({ nodeId: "c", type: NODE_TYPE.cnode, busy: true }),
    ]);
    expect(a.restartSite.available).toBe(false);
    expect(a.rf).toEqual({
      available: false,
      reason: "No amplifier node on this site",
    });
    expect(a.service).toEqual({
      available: false,
      reason: "No tower node on this site",
    });
  });
});
