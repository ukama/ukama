/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import {
  mapWithConcurrency,
  mapWithConcurrencySettled,
} from "../utils/concurrency";

const sleep = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

describe("mapWithConcurrency", () => {
  it("returns results in input order", async () => {
    const items = [30, 10, 20];
    const results = await mapWithConcurrency(items, async item => {
      await sleep(item);
      return item * 2;
    });
    expect(results).toEqual([60, 20, 40]);
  });

  it("never exceeds the concurrency limit", async () => {
    let inFlight = 0;
    let maxInFlight = 0;
    await mapWithConcurrency(
      Array.from({ length: 20 }, (_, i) => i),
      async () => {
        inFlight++;
        maxInFlight = Math.max(maxInFlight, inFlight);
        await sleep(5);
        inFlight--;
      },
      3
    );
    expect(maxInFlight).toBeLessThanOrEqual(3);
  });

  it("handles empty input", async () => {
    await expect(mapWithConcurrency([], async () => 1)).resolves.toEqual([]);
  });

  it("propagates rejections", async () => {
    await expect(
      mapWithConcurrency([1, 2], async item => {
        if (item === 2) throw new Error("boom");
        return item;
      })
    ).rejects.toThrow("boom");
  });
});

describe("mapWithConcurrencySettled", () => {
  it("returns per-item settled results without rejecting", async () => {
    const results = await mapWithConcurrencySettled([1, 2, 3], async item => {
      if (item === 2) throw new Error("boom");
      return item * 10;
    });
    expect(results[0]).toEqual({ status: "fulfilled", value: 10 });
    expect(results[1].status).toBe("rejected");
    expect(results[2]).toEqual({ status: "fulfilled", value: 30 });
  });
});
