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
