/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Verifies the assumption in docs/bff-screen-api-plan.md §4.2: Apollo's
 * RESTDataSource deduplicates identical concurrent GETs within one
 * (per-request) datasource instance, so composite sections sharing an
 * upstream call don't multiply upstream load. If an upgrade ever changes
 * this default, this test fails and the plan's caching strategy must be
 * revisited.
 */
import { BaseRESTDataSource } from "../datasource";

class TestDataSource extends BaseRESTDataSource {
  override baseURL = "http://upstream.test";

  fetchThing(): Promise<unknown> {
    return this.get("/v1/thing");
  }
}

describe("RESTDataSource GET deduplication", () => {
  const originalFetch = global.fetch;
  let fetchCount: number;

  beforeEach(() => {
    fetchCount = 0;
    global.fetch = jest.fn(async () => {
      fetchCount++;
      await new Promise(resolve => setTimeout(resolve, 10));
      return new Response(JSON.stringify({ ok: true }), {
        status: 200,
        headers: { "content-type": "application/json" },
      });
    }) as typeof fetch;
  });

  afterEach(() => {
    global.fetch = originalFetch;
  });

  it("collapses identical concurrent GETs into one upstream call", async () => {
    const ds = new TestDataSource();
    const [a, b, c] = await Promise.all([
      ds.fetchThing(),
      ds.fetchThing(),
      ds.fetchThing(),
    ]);
    expect(a).toEqual({ ok: true });
    expect(b).toEqual({ ok: true });
    expect(c).toEqual({ ok: true });
    expect(fetchCount).toBe(1);
  });

  it("does not dedupe across instances (per-request isolation)", async () => {
    const [r1, r2] = await Promise.all([
      new TestDataSource().fetchThing(),
      new TestDataSource().fetchThing(),
    ]);
    expect(r1).toEqual({ ok: true });
    expect(r2).toEqual({ ok: true });
    expect(fetchCount).toBe(2);
  });
});
