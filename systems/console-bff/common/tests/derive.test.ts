/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import "reflect-metadata";

import {
  countNodes,
  deriveSimPool,
  groupBy,
  groupNodeCountsBySite,
  subscriberActivity,
} from "../../dashboard/derive";
import { notImplementedSection, runSection } from "../../dashboard/section";
import { SectionErrorCode } from "../../dashboard/types";

const node = (connectivity: string, siteId?: string) =>
  ({
    status: { connectivity, state: "onboarded" },
    site: siteId ? { siteId } : undefined,
  }) as never;

describe("derive helpers", () => {
  it("countNodes counts online/offline by connectivity", () => {
    const counts = countNodes([
      node("Online"),
      node("online"),
      node("offline"),
    ]);
    expect(counts).toEqual({ total: 3, online: 2, offline: 1 });
  });

  it("groupNodeCountsBySite groups and counts per site", () => {
    const counts = groupNodeCountsBySite([
      node("online", "site-a"),
      node("offline", "site-a"),
      node("online", "site-b"),
      node("online"), // unassigned — excluded
    ]);
    expect(counts).toEqual([
      { siteId: "site-a", total: 2, online: 1, offline: 1 },
      { siteId: "site-b", total: 1, online: 1, offline: 0 },
    ]);
  });

  it("subscriberActivity derives active by attached sims", () => {
    expect(
      subscriberActivity([
        { sim: [{} as never] },
        { sim: [] },
        { sim: undefined },
      ])
    ).toEqual({ total: 3, active: 1, inactive: 2 });
  });

  it("deriveSimPool computes pctAssigned and lowStock", () => {
    expect(
      deriveSimPool({ total: 200, available: 5, consumed: 150 }, 10)
    ).toEqual({ pctAssigned: 75, lowStock: true });
    expect(deriveSimPool({ total: 0, available: 50, consumed: 0 }, 10)).toEqual(
      { pctAssigned: 0, lowStock: false }
    );
  });

  it("groupBy skips items without a key", () => {
    const groups = groupBy(
      [{ k: "a" }, { k: "a" }, { k: undefined }],
      item => item.k
    );
    expect(groups.get("a")).toHaveLength(2);
    expect(groups.size).toBe(1);
  });
});

describe("runSection / notImplementedSection", () => {
  it("returns value with null error on success", async () => {
    expect(await runSection("s", async () => 5)).toEqual({
      value: 5,
      error: null,
    });
  });

  it("returns typed error with null value on failure", async () => {
    const { value, error } = await runSection("kpis", async () => {
      throw { code: 404 };
    });
    expect(value).toBeNull();
    expect(error).toMatchObject({
      section: "kpis",
      code: SectionErrorCode.NOT_FOUND,
    });
  });

  it("notImplementedSection marks backend gaps", () => {
    expect(notImplementedSection("usage").error).toMatchObject({
      section: "usage",
      code: SectionErrorCode.NOT_IMPLEMENTED,
    });
  });
});
