/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Polled range metrics in the consolidated schema (plan Phase 4). This is
 * the query-only port of the subscriptions service's range fetch: same
 * upstream endpoint (`/v1/range/metrics/{key}`), same response mapping, but
 * NO WebSocket workers and NO GraphQL subscriptions — the console polls
 * (BUILD-PLAN §5.1·4). The standalone subscriptions service stays parked.
 */
import { Arg, Ctx, Field, InputType, Int, Query, Resolver } from "type-graphql";

import { STATS_TYPE } from "../../common/enums";
import { mapWithConcurrency } from "../../common/utils/concurrency";
import type { AppContext } from "../../server/context";
import { getNodeMetricRange } from "../../subscriptions/datasource/subscriptions-api";
import type {
  GetMetricsStatInput,
  MetricRes,
  MetricsRes,
} from "../../subscriptions/resolvers/types";
import { MetricsRes as MetricsResType } from "../../subscriptions/resolvers/types";
import { ServiceUrlResolver } from "../baseUrls";
import { isMockKey, metricMeta } from "../metrics/catalog";
import { mockRangeValues } from "../metrics/mock";

/** True when the node isn't online (offline/unknown) — no telemetry expected. */
const isNodeOffline = async (
  ctx: AppContext,
  nodeId: string
): Promise<boolean> => {
  try {
    const urls = new ServiceUrlResolver(ctx.headers.orgName);
    const url = await urls.url("node");
    const node = await ctx.dataSources.node.getNode(url, { id: nodeId });
    return node.status?.connectivity?.toLowerCase() !== "online";
  } catch {
    return false;
  }
};

/** Normalize raw units to console-friendly ones (applied to every series so
 *  it's consistent everywhere). bps → Mbps: scale values + threshold by 1e6,
 *  preserving the -1 gap-fill sentinel. */
const UNIT_SCALE: Record<string, { unit: string; div: number }> = {
  bps: { unit: "Mbps", div: 1e6 },
};
const normalizeUnits = (m: MetricRes): MetricRes => {
  const rule = m.unit ? UNIT_SCALE[m.unit] : undefined;
  if (!rule) return m;
  const scale = (v: number) => (v === -1 ? -1 : v / rule.div);
  return {
    ...m,
    unit: rule.unit,
    values: (m.values ?? []).map(([t, v]) => [t, scale(v)]),
    threshold: m.threshold
      ? {
          min: m.threshold.min / rule.div,
          normal: m.threshold.normal / rule.div,
          max: m.threshold.max / rule.div,
        }
      : m.threshold,
  };
};

/** Backfill the presentation metadata the upstream omits (label always; unit/
 *  format/threshold only when missing) so the console renders from data, then
 *  normalize units. */
const enrich = (m: MetricRes): MetricRes =>
  normalizeUnits({
    ...m,
    label: m.label ?? (metricMeta(m.type).label || m.type),
    unit: m.unit ?? metricMeta(m.type).unit,
    format: m.format ?? metricMeta(m.type).format,
    threshold:
      m.threshold ??
      (metricMeta(m.type).threshold
        ? { ...metricMeta(m.type).threshold! }
        : undefined),
  });

const MAX_KEYS = 10;

/**
 * Resolution control for range charts. The console's Day/Week/Month filter
 * sets the from/to window; we derive a Prometheus `step` from the span so the
 * series is a bounded ~TARGET_POINTS samples regardless of range, instead of
 * the old hardcoded 1s (which pulled 86 400 points for a single Day).
 *   Day   (86 400s)  → step 720s  (12 min)  → ~120 points
 *   Week  (604 800s) → step 5 040s (84 min) → ~120 points
 *   Month (2 592 000s) → step 21 600s (6 h) → ~120 points
 * MIN_STEP keeps us at/above Prometheus' 15s scrape interval so we never
 * oversample. MAX_STEP caps very wide windows.
 */
const TARGET_POINTS = 120;
const MIN_STEP = 60;
const MAX_STEP = 86_400;

/** Bucketing step (seconds) for a [from,to] window, clamped to sane bounds. */
const deriveStep = (from: number, to: number): number => {
  const span = Math.max(0, to - from);
  const step = Math.ceil(span / TARGET_POINTS);
  return Math.min(MAX_STEP, Math.max(MIN_STEP, step));
};

@InputType()
export class MetricsRangeInput {
  /** Metric keys, e.g. ["uptime", "cpu_temperature"]. Max 10 per request. */
  @Field(() => [String])
  keys: string[];

  /** Epoch seconds (must be > 0). */
  @Field(() => Int)
  from: number;

  @Field(() => Int, { nullable: true })
  to?: number;

  @Field({ nullable: true })
  nodeId?: string;

  /** Prometheus aggregation, default "avg". */
  @Field({ nullable: true })
  operation?: string;

  /** Optional resolution override (seconds). Omit to auto-derive from the
   *  [from,to] window — the console's Day/Week/Month filter drives that. */
  @Field(() => Int, { nullable: true })
  step?: number;
}

@Resolver()
export class MetricsRangeResolver {
  @Query(() => MetricsResType)
  async metricsRange(
    @Arg("data") data: MetricsRangeInput,
    @Ctx() ctx: AppContext
  ): Promise<MetricsRes> {
    if (data.from <= 0) {
      throw new Error("Argument 'from' must be a positive epoch timestamp.");
    }
    const keys = data.keys.slice(0, MAX_KEYS);
    const to = data.to ?? Math.floor(Date.now() / 1000);
    const step =
      data.step && data.step > 0 ? data.step : deriveStep(data.from, to);
    const scope = data.nodeId ?? ctx.headers.orgName;

    // An offline/unknown node reports no telemetry — return empty series so
    // the console shows a "no data" state instead of a fabricated chart.
    const nodeOffline = data.nodeId
      ? await isNodeOffline(ctx, data.nodeId)
      : false;

    const liveKeys = keys.filter(k => !isMockKey(k));
    const mockKeys = keys.filter(k => isMockKey(k));

    // Mocked keys: synthesize MetricRes directly (no upstream call).
    const mocked: MetricRes[] = mockKeys.map(key =>
      enrich({
        success: true,
        msg: "mock",
        type: key,
        nodeId: data.nodeId,
        values: nodeOffline
          ? []
          : mockRangeValues(key, scope, data.from, to, step),
      } as MetricRes)
    );

    // Live keys: real metric service, then backfill presentation metadata.
    let live: MetricRes[] = [];
    if (liveKeys.length > 0) {
      const urls = new ServiceUrlResolver(ctx.headers.orgName);
      const baseURL = await urls.url("metrics");
      const args = {
        from: data.from,
        to,
        step,
        nodeId: data.nodeId,
        operation: data.operation ?? "avg",
        orgName: ctx.headers.orgName,
        userId: ctx.headers.userId,
        type: STATS_TYPE.HOME,
        withSubscription: false,
      } as GetMetricsStatInput;
      const results = await mapWithConcurrency(liveKeys, key =>
        getNodeMetricRange(baseURL, key, args)
      );
      // Always emit one entry per requested key: when the upstream returns no
      // series (e.g. the gateway has no mapping/data for it yet), fall back to
      // an empty, enriched placeholder so the console still gets the catalog
      // label/unit (no raw key shown) and renders a "no data" state.
      live = liveKeys.flatMap((key, i) => {
        const found = results[i]?.metrics ?? [];
        if (found.length > 0) return found.map(enrich);
        return [
          enrich({
            success: false,
            msg: "no-data",
            type: key,
            nodeId: data.nodeId,
            values: [],
          } as MetricRes),
        ];
      });
    }

    // One series per (key, node/site). The gateway's `avg(...) without(job,
    // instance,receive,tenant_id)` keeps stray labels (e.g. the dummy's
    // `metric` label, exported_* from multi-path ingestion), so a single key
    // can come back as multiple groups that we then stamp with the same
    // nodeId — collapse them so each chart key yields one line.
    const seen = new Set<string>();
    const metrics = [...mocked, ...live].filter(m => {
      const k = `${m.type}|${m.nodeId ?? ""}|${m.siteId ?? ""}`;
      if (seen.has(k)) return false;
      seen.add(k);
      return true;
    });

    return { metrics };
  }
}
