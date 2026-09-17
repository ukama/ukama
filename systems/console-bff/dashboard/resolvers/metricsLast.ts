/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Latest value per metric key for one node, read through the metrics
 * gateway's instant endpoint (`/v1/last/metrics/{key}?node=`). Pairs with
 * metricsRange: that one draws charts over a window, this one answers "what
 * is it right now" for rail values that have no chart, such as uptime.
 */
import { Arg, Ctx, Field, InputType, Query, Resolver } from "type-graphql";

import { logger } from "../../common/logger";
import { mapWithConcurrency } from "../../common/utils/concurrency";
import type { AppContext } from "../../server/context";
import { ServiceUrlResolver } from "../baseUrls";
import { metricMeta } from "../metrics/catalog";
import { runSection } from "../section";
import { KpiEntryDto, KpisSection } from "./types";

const MAX_KEYS = 10;

@InputType()
export class MetricsLastInput {
  /** Metric keys, e.g. ["com_uptime"]. Max 10 per request. */
  @Field(() => [String])
  keys: string[];

  @Field()
  nodeId: string;
}

@Resolver()
export class MetricsLastResolver {
  @Query(() => KpisSection)
  async metricsLast(
    @Arg("data") data: MetricsLastInput,
    @Ctx() ctx: AppContext
  ): Promise<KpisSection> {
    const keys = data.keys.slice(0, MAX_KEYS);
    const { value, error } = await runSection("metricsLast", async () => {
      const urls = new ServiceUrlResolver(ctx.headers.orgName);
      // The pushgateway keeps serving a dead node's last sample, so an
      // offline node reads as no data rather than a stale value.
      const nodeUrl = await urls.url("node");
      const node = await ctx.dataSources.node.getNode(nodeUrl, {
        id: data.nodeId,
      });
      const online = node.status?.connectivity?.toLowerCase() === "online";
      const metricsUrl = await urls.url("metrics");

      const results = await mapWithConcurrency(keys, async key => {
        if (!online) {
          return {
            type: key,
            value: [0, 0] as [number, number],
            success: false,
          };
        }
        try {
          return await ctx.dataSources.metric.getNodeLast(
            metricsUrl,
            key,
            data.nodeId
          );
        } catch (e) {
          logger.warn(`[metricsLast] '${key}' failed: ${e}`);
          return {
            type: key,
            value: [0, 0] as [number, number],
            success: false,
          };
        }
      });

      return results.map((r): KpiEntryDto => {
        const meta = metricMeta(r.type);
        return {
          key: r.type,
          timestamp: r.value[0],
          value: r.value[1],
          success: r.success,
          label: meta.label || r.type,
          unit: meta.unit,
          format: meta.format,
          threshold: meta.threshold ?? null,
        };
      });
    });
    return { metrics: value, error };
  }
}
