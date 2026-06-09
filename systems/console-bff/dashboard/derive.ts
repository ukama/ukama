/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Pure derivation helpers for composite sections (plan §2.6: UI-ready DTOs —
 * the console renders, it does not compute). Kept dependency-free so they
 * are trivially unit-testable.
 */
import { Node } from "../node/resolvers/types";
import { SimPoolStatsDto } from "../sim/resolver/types";
import { SubscriberDto } from "../subscriber/resolver/types";

export const SIM_POOL_LOW_STOCK_THRESHOLD = parseInt(
  process.env.SIM_POOL_LOW_STOCK_THRESHOLD ?? "10"
);

export interface NodeCounts {
  total: number;
  online: number;
  offline: number;
}

export const countNodes = (nodes: readonly Node[]): NodeCounts => {
  let online = 0;
  for (const node of nodes) {
    if (node.status?.connectivity?.toLowerCase() === "online") online++;
  }
  return { total: nodes.length, online, offline: nodes.length - online };
};

export interface SiteNodeCount extends NodeCounts {
  siteId: string;
}

export const groupNodeCountsBySite = (
  nodes: readonly Node[]
): SiteNodeCount[] => {
  const bySite = new Map<string, Node[]>();
  for (const node of nodes) {
    const siteId = node.site?.siteId;
    if (!siteId) continue;
    const group = bySite.get(siteId);
    if (group) group.push(node);
    else bySite.set(siteId, [node]);
  }
  return Array.from(bySite.entries(), ([siteId, siteNodes]) => ({
    siteId,
    ...countNodes(siteNodes),
  }));
};

export interface SubscriberActivity {
  total: number;
  active: number;
  inactive: number;
}

/** A subscriber is "active" when at least one SIM is attached. */
export const subscriberActivity = (
  subscribers: readonly Pick<SubscriberDto, "sim">[]
): SubscriberActivity => {
  let active = 0;
  for (const sub of subscribers) {
    if ((sub.sim?.length ?? 0) > 0) active++;
  }
  return {
    total: subscribers.length,
    active,
    inactive: subscribers.length - active,
  };
};

export interface SimPoolDerived {
  pctAssigned: number;
  lowStock: boolean;
}

export const deriveSimPool = (
  stats: Pick<SimPoolStatsDto, "total" | "available" | "consumed">,
  lowStockThreshold: number = SIM_POOL_LOW_STOCK_THRESHOLD
): SimPoolDerived => ({
  pctAssigned:
    stats.total > 0 ? Math.round((stats.consumed / stats.total) * 100) : 0,
  lowStock: stats.available < lowStockThreshold,
});

/* ------------------------- commerce (Phase 3) ---------------------------- */

export interface PaymentLike {
  amount: string;
  status: string;
  paidAt: string;
  itemType: string;
  itemId: string;
}

const PAID_STATUSES = new Set(["success", "paid", "completed", "processed"]);

export const isPaidPayment = (payment: Pick<PaymentLike, "status">): boolean =>
  PAID_STATUSES.has(payment.status.toLowerCase());

const monthKey = (iso: string): string => iso.slice(0, 7); // YYYY-MM

export interface RevenueSummary {
  totalPaid: number;
  totalPending: number;
  monthPaid: number;
  prevMonthPaid: number;
  /** Month-over-month change of paid revenue, percent (null if no base). */
  momPct: number | null;
}

export const summarizeRevenue = (
  payments: readonly PaymentLike[],
  now: Date = new Date()
): RevenueSummary => {
  const thisMonth = now.toISOString().slice(0, 7);
  const prev = new Date(Date.UTC(now.getUTCFullYear(), now.getUTCMonth() - 1));
  const prevMonth = prev.toISOString().slice(0, 7);
  let totalPaid = 0;
  let totalPending = 0;
  let monthPaid = 0;
  let prevMonthPaid = 0;
  for (const payment of payments) {
    const amount = parseFloat(payment.amount) || 0;
    if (isPaidPayment(payment)) {
      totalPaid += amount;
      const key = monthKey(payment.paidAt ?? "");
      if (key === thisMonth) monthPaid += amount;
      else if (key === prevMonth) prevMonthPaid += amount;
    } else {
      totalPending += amount;
    }
  }
  const momPct =
    prevMonthPaid > 0
      ? Math.round(((monthPaid - prevMonthPaid) / prevMonthPaid) * 100)
      : null;
  return { totalPaid, totalPending, monthPaid, prevMonthPaid, momPct };
};

export interface PlanLike {
  uuid: string;
  name: string;
  active: boolean;
  amount: number;
  currency: string;
}

export interface SimPackageLike {
  packageId?: string;
  isActive?: boolean;
}

export interface PlanStats {
  packageId: string;
  name: string;
  amount: number;
  currency: string;
  active: boolean;
  /** Active attachments (null when no network scope was given). */
  attachCount: number | null;
  /** Paid revenue attributed to this plan (payments with itemType=package). */
  revenue: number;
  revenueSharePct: number;
}

export const derivePlanStats = (
  plans: readonly PlanLike[],
  payments: readonly PaymentLike[],
  activeSimPackages: readonly SimPackageLike[] | null
): PlanStats[] => {
  const revenueByPlan = new Map<string, number>();
  let totalRevenue = 0;
  for (const payment of payments) {
    if (payment.itemType?.toLowerCase() !== "package") continue;
    if (!isPaidPayment(payment)) continue;
    const amount = parseFloat(payment.amount) || 0;
    revenueByPlan.set(
      payment.itemId,
      (revenueByPlan.get(payment.itemId) ?? 0) + amount
    );
    totalRevenue += amount;
  }
  const attachByPlan = activeSimPackages
    ? groupBy(
        activeSimPackages.filter(p => p.isActive),
        p => p.packageId
      )
    : null;
  return plans.map(plan => {
    const revenue = revenueByPlan.get(plan.uuid) ?? 0;
    return {
      packageId: plan.uuid,
      name: plan.name,
      amount: plan.amount,
      currency: plan.currency,
      active: plan.active,
      attachCount: attachByPlan
        ? (attachByPlan.get(plan.uuid)?.length ?? 0)
        : null,
      revenue,
      revenueSharePct:
        totalRevenue > 0 ? Math.round((revenue / totalRevenue) * 100) : 0,
    };
  });
};

export interface CategoryCount {
  category: string;
  count: number;
}

export const countByCategory = (
  items: readonly { category?: string }[]
): CategoryCount[] =>
  Array.from(
    groupBy(items, item => item.category || "uncategorized").entries(),
    ([category, group]) => ({ category, count: group.length })
  );

/** O(n+m) group-join of sims onto their owner, keyed by `subscriberId`. */
export const groupBy = <T>(
  items: readonly T[],
  key: (item: T) => string | undefined
): Map<string, T[]> => {
  const groups = new Map<string, T[]>();
  for (const item of items) {
    const k = key(item);
    if (!k) continue;
    const group = groups.get(k);
    if (group) group.push(item);
    else groups.set(k, [item]);
  }
  return groups;
};
