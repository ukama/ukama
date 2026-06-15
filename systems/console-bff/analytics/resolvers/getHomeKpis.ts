/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { HomeKpis, HomeViewInput } from "./types/home";

@Resolver()
export class GetHomeKpisResolver {
  @Query(() => HomeKpis)
  async getHomeKpis(
    @Arg("data") data: HomeViewInput,
    @Ctx() ctx: AppContext
  ): Promise<HomeKpis> {
    const baseURL = await ctx.urls.url("analytics");
    const result = await ctx.dataSources.analytics.getHomeKpis(baseURL, data);

    // "Active customers" = the real active-subscriber count from the metric
    // service (system-scoped `subscribers_active`). Same source for both Home
    // lenses; overrides the analytics value so the figure is consistent.
    try {
      const metricsURL = await ctx.urls.url("metrics");
      const m = await ctx.dataSources.metric.getLatestMetric(
        metricsURL,
        "subscribers_active"
      );
      if (m.success) {
        const value = Math.round(Number(m.value?.[1] ?? 0));
        const kpis = [...result.kpis];
        const i = kpis.findIndex(k => k.key === "active_customers");
        if (i >= 0) {
          kpis[i] = { ...kpis[i], value, formatted: String(value), stale: false };
        } else {
          kpis.push({ key: "active_customers", value, formatted: String(value) });
        }
        return { ...result, kpis };
      }
    } catch {
      // Metric service unavailable — keep the analytics value.
    }
    return result;
  }
}
