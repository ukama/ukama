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

/** Backfill the presentation metadata the upstream omits (label always; unit/
 *  format/threshold only when missing) so the console renders from data. */
const enrich = (m: MetricRes): MetricRes => {
  const meta = metricMeta(m.type);
  return {
    ...m,
    label: m.label ?? (meta.label || m.type),
    unit: m.unit ?? meta.unit,
    format: m.format ?? meta.format,
    threshold:
      m.threshold ?? (meta.threshold ? { ...meta.threshold } : undefined),
  };
};

const MAX_KEYS = 10;

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
    const step = Math.max(60, Math.floor((to - data.from) / 48));
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
      live = results.flatMap(res => res.metrics).map(enrich);
    }

    return { metrics: [...mocked, ...live] };
  }
}
