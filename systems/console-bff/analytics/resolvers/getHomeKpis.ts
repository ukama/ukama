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

    const kpis = [...result.kpis];
    const upsert = (key: string, value: number): void => {
      const i = kpis.findIndex(k => k.key === key);
      if (i >= 0) {
        kpis[i] = { ...kpis[i], value, formatted: String(value), stale: false };
      } else {
        kpis.push({ key, value, formatted: String(value) });
      }
    };

    // "Active customers" = the real active-subscriber count from the metric
    // service (system-scoped `subscribers_active`). Overrides the analytics
    // value so the figure is consistent across Home lenses.
    try {
      const metricsURL = await ctx.urls.url("metrics");
      const m = await ctx.dataSources.metric.getLatestMetric(
        metricsURL,
        "subscribers_active"
      );
      if (m.success)
        upsert("active_customers", Math.round(Number(m.value?.[1] ?? 0)));
    } catch {
      // Metric service unavailable — keep the analytics value.
    }

    // "Total customers" = the real all-states subscriber count from the
    // subscriber service (analytics has no total source). Needs a networkId;
    // an org-wide business home without one leaves the key absent → the console
    // shows "—" rather than a fabricated figure.
    if (data.networkId) {
      try {
        const subURL = await ctx.urls.url("subscriber");
        const subs = await ctx.dataSources.subscriber.getSubscribersByNetwork(
          subURL,
          data.networkId
        );
        upsert("customers_total", subs.subscribers.length);
      } catch {
        // Subscriber service unavailable — leave customers_total absent.
      }
    }

    return { ...result, kpis };
  }
}
