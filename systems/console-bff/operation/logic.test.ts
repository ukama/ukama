/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { NODE_CONNECTIVITY, NODE_TYPE } from "../common/enums";
import {
  SiteNodeStatus,
  buildSiteActions,
  isLockBusy,
  leaseExpired,
  nodeResourceKey,
  toNodeStatus,
} from "./logic";
import {
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
  over: Partial<SiteNodeStatus>
): SiteNodeStatus => ({
  nodeId: "n",
  busy: false,
  connectivity: NODE_CONNECTIVITY.Online,
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
      connectivity: undefined,
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

describe("buildSiteActions — independent async release", () => {
  it("keeps service unavailable after lease expiry until the manager releases it", () => {
    const expired = op({ leaseExpiresAt: new Date(NOW - 1000).toISOString() });
    const tower = toNodeStatus(
      { id: "t", type: NODE_TYPE.tnode, lock: lock({ operation: expired }), connectivity: NODE_CONNECTIVITY.Online },
      NOW
    );
    const amp = status({ nodeId: "a", type: NODE_TYPE.anode });
    const waiting = buildSiteActions([tower, amp]);

    expect(tower.operation?.id).toBe(expired.id);
    expect(waiting.service.available).toBe(false);
    expect(waiting.restartSite.available).toBe(false);
    expect(waiting.rf.available).toBe(true);

    const released = toNodeStatus(
      { id: "t", type: NODE_TYPE.tnode, lock: { locked: false }, connectivity: NODE_CONNECTIVITY.Online },
      NOW
    );
    const ready = buildSiteActions([released, amp]);
    expect(released.operation).toBe(undefined);
    expect(ready.service.available).toBe(true);
    expect(ready.restartSite.available).toBe(true);
  });

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
    expect(a.service.available).toBe(true);
  });

  it("tower busy → service and restartSite locked, RF stays available", () => {
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
    expect(a.rf.available).toBe(true);
    expect(a.restartSite.available).toBe(false);
  });

  it("missing role node → action unavailable with a clear reason", () => {
    const a = buildSiteActions([
      status({ nodeId: "t", type: NODE_TYPE.tnode }),
    ]);
    expect(a.rf.available).toBe(false);
    expect(a.rf.reason).toBe("No amplifier node on this site");
  });
});

describe("buildSiteActions — connectivity and controller dependencies", () => {
  const site = (): SiteNodeStatus[] => [
    status({ nodeId: "t", type: NODE_TYPE.tnode }),
    status({ nodeId: "a", type: NODE_TYPE.anode }),
    status({ nodeId: "c", type: NODE_TYPE.cnode }),
  ];

  const expectBlocked = (actions: ReturnType<typeof buildSiteActions>) => {
    for (const action of [actions.restartSite, actions.rf, actions.service]) {
      expect(action.available).toBe(false);
      expect(action.reason?.trim().length).toBeGreaterThan(0);
    }
  };

  it("all offline and unlocked disables every action with a reason", () => {
    const nodes = site().map(n => ({
      ...n, connectivity: NODE_CONNECTIVITY.Offline,
    }));
    expectBlocked(buildSiteActions(nodes));
  });

  it("offline amplifier blocks RF and restart, but leaves service available", () => {
    const nodes = site();
    nodes[1].connectivity = NODE_CONNECTIVITY.Offline;
    const actions = buildSiteActions(nodes);
    expect(actions.restartSite.available).toBe(false);
    expect(actions.rf.available).toBe(false);
    expect(actions.rf.reason).toContain("Node a ");
    expect(actions.service.available).toBe(true);
  });

  it("offline tower blocks service and restart, but leaves RF available", () => {
    const nodes = site();
    nodes[0].connectivity = NODE_CONNECTIVITY.Offline;
    const actions = buildSiteActions(nodes);
    expect(actions.restartSite.available).toBe(false);
    expect(actions.service.available).toBe(false);
    expect(actions.service.reason).toContain("Node t ");
    expect(actions.rf.available).toBe(true);
  });

  it("offline controller disables every action with a reason", () => {
    const nodes = site();
    nodes[2].connectivity = NODE_CONNECTIVITY.Offline;
    expectBlocked(buildSiteActions(nodes));
  });

  it("unknown or missing connectivity never enables an action", () => {
    for (const connectivity of [NODE_CONNECTIVITY.Unknown, undefined, ""]) {
      expectBlocked(buildSiteActions(site().map(n => ({ ...n, connectivity }))));
    }
  });

  it("busy controller disables all actions until its lock clears", () => {
    const nodes = site();
    nodes[2].busy = true;
    nodes[2].operation = op({ type: "RestartNode", requestedBy: "sam" });
    const busy = buildSiteActions(nodes);
    expectBlocked(busy);
    expect(busy.rf.reason).toContain("sam");
    expect(busy.service.reason).toContain("sam");

    nodes[2].busy = false;
    nodes[2].operation = undefined;
    const ready = buildSiteActions(nodes);
    expect(ready).toEqual({
      restartSite: { available: true },
      rf: { available: true },
      service: { available: true },
    });
  });

  it("clearing a lock while still offline does not enable actions", () => {
    const nodes = site();
    nodes[2] = toNodeStatus({
      id: "c", type: NODE_TYPE.cnode,
      connectivity: NODE_CONNECTIVITY.Offline,
      lock: { locked: false },
    }, NOW);
    expectBlocked(buildSiteActions(nodes));
    nodes[2].connectivity = NODE_CONNECTIVITY.Online;
    const ready = buildSiteActions(nodes);
    expect(ready.restartSite.available).toBe(true);
    expect(ready.rf.available).toBe(true);
    expect(ready.service.available).toBe(true);
  });

  it("failed lock reads preserve connectivity and existing fail-open behavior", () => {
    const nodes = site().map(n => toNodeStatus({
      id: n.nodeId, type: n.type,
      connectivity: NODE_CONNECTIVITY.Offline, failed: true,
    }, NOW));
    expectBlocked(buildSiteActions(nodes));
    const ready = buildSiteActions(nodes.map(n => ({
      ...n, connectivity: NODE_CONNECTIVITY.Online,
    })));
    expect(ready.restartSite.available).toBe(true);
    expect(ready.rf.available).toBe(true);
    expect(ready.service.available).toBe(true);
  });

  it("an empty node list disables every action with a reason", () => {
    expectBlocked(buildSiteActions([]));
  });

  it("checks every controller, even if an earlier controller is idle", () => {
    const nodes = site();
    nodes.push(status({
      nodeId: "c2", type: NODE_TYPE.cnode, busy: true, operation: op(),
    }));
    expectBlocked(buildSiteActions(nodes));
  });
});
