/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import "reflect-metadata";

import {
  countByCategory,
  countNodes,
  derivePlanStats,
  deriveSimPool,
  groupBy,
  groupNodeCountsBySite,
  subscriberActivity,
  summarizeRevenue,
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

describe("commerce derivations (Phase 3)", () => {
  const pay = (
    amount: string,
    status: string,
    paidAt: string,
    itemType = "invoice",
    itemId = "x"
  ) => ({ amount, status, paidAt, itemType, itemId });

  it("summarizeRevenue splits paid/pending and computes MoM", () => {
    const now = new Date("2026-06-04T00:00:00Z");
    const summary = summarizeRevenue(
      [
        pay("100", "success", "2026-06-01"),
        pay("50", "paid", "2026-05-20"),
        pay("25", "pending", "2026-06-02"),
        pay("10", "completed", "2026-04-01"),
      ],
      now
    );
    expect(summary).toEqual({
      totalPaid: 160,
      totalPending: 25,
      monthPaid: 100,
      prevMonthPaid: 50,
      momPct: 100,
    });
  });

  it("summarizeRevenue returns null MoM without a base month", () => {
    const now = new Date("2026-06-04T00:00:00Z");
    expect(
      summarizeRevenue([pay("100", "success", "2026-06-01")], now).momPct
    ).toBeNull();
  });

  it("derivePlanStats attributes package revenue and attach counts", () => {
    const plans = [
      {
        uuid: "p1",
        name: "Standard",
        active: true,
        amount: 10,
        currency: "USD",
      },
      { uuid: "p2", name: "Starter", active: true, amount: 5, currency: "USD" },
    ];
    const payments = [
      pay("30", "success", "2026-06-01", "package", "p1"),
      pay("10", "success", "2026-06-01", "package", "p2"),
      pay("99", "pending", "2026-06-01", "package", "p1"), // unpaid: excluded
      pay("70", "success", "2026-06-01", "invoice", "p1"), // not a package
    ];
    const sims = [
      { packageId: "p1", isActive: true },
      { packageId: "p1", isActive: true },
      { packageId: "p2", isActive: false },
    ];
    const stats = derivePlanStats(plans, payments, sims);
    expect(stats[0]).toMatchObject({
      packageId: "p1",
      attachCount: 2,
      revenue: 30,
      revenueSharePct: 75,
    });
    expect(stats[1]).toMatchObject({
      packageId: "p2",
      attachCount: 0,
      revenue: 10,
      revenueSharePct: 25,
    });
  });

  it("derivePlanStats keeps attachCount null without network scope", () => {
    const stats = derivePlanStats(
      [{ uuid: "p1", name: "S", active: true, amount: 1, currency: "USD" }],
      [],
      null
    );
    expect(stats[0].attachCount).toBeNull();
    expect(stats[0].revenueSharePct).toBe(0);
  });

  it("countByCategory groups components", () => {
    expect(
      countByCategory([
        { category: "ACCESS" },
        { category: "ACCESS" },
        { category: undefined },
      ])
    ).toEqual([
      { category: "ACCESS", count: 2 },
      { category: "uncategorized", count: 1 },
    ]);
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
