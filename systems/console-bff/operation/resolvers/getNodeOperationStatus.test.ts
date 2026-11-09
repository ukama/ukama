/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import "reflect-metadata";

import type { AppContext } from "../../server/context";
import { GetNodeOperationStatusResolver } from "./getNodeOperationStatus";
import { OperationDto, ResourceLockDto } from "./types";

const op = (resourceKey: string, requestedBy: string): OperationDto => ({
  id: `op-${resourceKey}`,
  type: "UpdateSoftware",
  system: "node",
  status: "running",
  fencingToken: 1,
  resourceKey,
  requestedBy,
});

interface Fake {
  locks?: Record<string, ResourceLockDto | Error>;
  siteId?: string | Error;
  siteNodes?: string[] | Error;
}

const fakeContext = (fake: Fake) => {
  const calls = { getNode: 0, getNodes: 0, locks: [] as string[] };
  const result = <T>(v: T | Error | undefined): Promise<T | undefined> =>
    v instanceof Error ? Promise.reject(v) : Promise.resolve(v);

  const ctx = {
    urls: { url: async (name: string) => `http://${name}` },
    dataSources: {
      operation: {
        getResourceLock: (_url: string, key: string) => {
          calls.locks.push(key);
          return result(fake.locks?.[key] ?? { locked: false });
        },
      },
      node: {
        getNode: (_url: string, { id }: { id: string }) => {
          calls.getNode++;
          return result(fake.siteId).then(siteId => ({
            id,
            site: { siteId: siteId ?? "" },
          }));
        },
        getNodes: () => {
          calls.getNodes++;
          return result(fake.siteNodes).then(ids => ({
            nodes: (ids ?? []).map(id => ({ id })),
          }));
        },
      },
    },
  } as unknown as AppContext;

  return { ctx, calls };
};

const query = (ctx: AppContext) =>
  new GetNodeOperationStatusResolver().getNodeOperationStatus(
    { nodeId: "t" },
    ctx
  );

describe("getNodeOperationStatus", () => {
  it("busy node is blocked by its own lock without reading its site", async () => {
    const own = op("node:t", "own");
    const { ctx, calls } = fakeContext({
      locks: { "node:t": { locked: true, operation: own } },
    });

    const s = await query(ctx);

    expect(s).toEqual({
      nodeId: "t",
      busy: true,
      operation: own,
      blocked: true,
      blockingOperation: own,
    });
    expect(calls.getNode).toBe(0);
  });

  it("idle node is blocked by a busy sibling", async () => {
    const update = op("node:c", "sam");
    const { ctx, calls } = fakeContext({
      siteId: "s1",
      siteNodes: ["t", "a", "c"],
      locks: { "node:c": { locked: true, operation: update } },
    });

    const s = await query(ctx);

    expect(s.busy).toBe(false);
    expect(s.operation).toBe(undefined);
    expect(s.blocked).toBe(true);
    expect(s.blockingOperation).toEqual(update);
    expect(calls.locks.sort()).toEqual(["node:a", "node:c", "node:t"]);
  });

  it("idle site leaves the node unblocked", async () => {
    const { ctx } = fakeContext({ siteId: "s1", siteNodes: ["t", "a"] });

    const s = await query(ctx);

    expect(s.blocked).toBe(false);
    expect(s.blockingOperation).toBe(undefined);
  });

  it("node without a site is not checked against siblings", async () => {
    const { ctx, calls } = fakeContext({ siteId: "" });

    const s = await query(ctx);

    expect(s.blocked).toBe(false);
    expect(calls.getNodes).toBe(0);
  });

  it("failed reads fail open", async () => {
    const fakes: Fake[] = [
      { locks: { "node:t": new Error("down") }, siteId: new Error("down") },
      { siteId: "s1", siteNodes: new Error("down") },
      {
        siteId: "s1",
        siteNodes: ["t", "c"],
        locks: { "node:c": new Error("down") },
      },
    ];
    for (const fake of fakes) {
      const s = await query(fakeContext(fake).ctx);
      expect(s.busy).toBe(false);
      expect(s.blocked).toBe(false);
    }
  });
});
