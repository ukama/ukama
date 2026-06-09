/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import {
  Arg,
  Ctx,
  FieldResolver,
  Int,
  Query,
  Resolver,
  Root,
} from "type-graphql";

import type { AppContext } from "../../server/context";
import { ServiceUrlResolver } from "../baseUrls";
import { countNodes, subscriberActivity } from "../derive";
import { NETWORK_KPI_KEYS, fetchLatestKpis } from "../kpis";
import { runSection } from "../section";
import {
  AlertsSection,
  KpisSection,
  NetworkOverview,
  NetworkSection,
  NodeStatsSection,
  SitesSection,
  SubscriberStatsSection,
} from "./types";

const MAX_ALERTS = 20;

type OverviewRoot = NetworkOverview & { _urls: ServiceUrlResolver };

/**
 * Network-lens home composite (plan §3.1). Core = networkId only; every
 * section is lazy and pays only when selected. Serves: Network home,
 * Business home (network section), Biz network screen.
 */
@Resolver(() => NetworkOverview)
export class NetworkOverviewResolver {
  @Query(() => NetworkOverview)
  networkOverview(
    @Arg("networkId") networkId: string,
    @Ctx() ctx: AppContext
  ): NetworkOverview {
    const root: OverviewRoot = Object.assign(new NetworkOverview(), {
      networkId,
      _urls: new ServiceUrlResolver(ctx.headers.orgName),
    });
    return root;
  }

  @FieldResolver(() => NetworkSection)
  async network(
    @Root() root: OverviewRoot,
    @Ctx() ctx: AppContext
  ): Promise<NetworkSection> {
    const { value, error } = await runSection("network", async () => {
      const url = await root._urls.url("network");
      return ctx.dataSources.network.getNetwork(url, root.networkId);
    });
    return { network: value, error };
  }

  @FieldResolver(() => NodeStatsSection)
  async nodeStats(
    @Root() root: OverviewRoot,
    @Ctx() ctx: AppContext
  ): Promise<NodeStatsSection> {
    const { value, error } = await runSection("nodeStats", async () => {
      const url = await root._urls.url("node");
      const res = await ctx.dataSources.node.getNodesByNetwork(
        url,
        root.networkId
      );
      return countNodes(res.nodes);
    });
    return { ...value, error };
  }

  @FieldResolver(() => SitesSection)
  async siteStats(
    @Root() root: OverviewRoot,
    @Ctx() ctx: AppContext
  ): Promise<SitesSection> {
    const { value, error } = await runSection("siteStats", async () => {
      const url = await root._urls.url("site");
      const res = await ctx.dataSources.site.getSites(url, {
        networkId: root.networkId,
      });
      return res.sites;
    });
    return { sites: value, error };
  }

  @FieldResolver(() => SubscriberStatsSection)
  async subscriberStats(
    @Root() root: OverviewRoot,
    @Ctx() ctx: AppContext
  ): Promise<SubscriberStatsSection> {
    const { value, error } = await runSection("subscriberStats", async () => {
      const url = await root._urls.url("subscriber");
      const [sims, subs] = await Promise.all([
        ctx.dataSources.subscriber.getSimsByNetwork(url, root.networkId),
        ctx.dataSources.subscriber.getSubscribersByNetwork(url, root.networkId),
      ]);
      const withSims = subs.subscribers.map(sub => ({
        sim: sims.sims.filter(sim => sim.subscriberId === sub.uuid),
      }));
      return subscriberActivity(withSims);
    });
    return { ...value, error };
  }

  @FieldResolver(() => AlertsSection)
  async latestAlerts(
    @Root() root: OverviewRoot,
    @Ctx() ctx: AppContext,
    @Arg("limit", () => Int, { defaultValue: 5 }) limit: number
  ): Promise<AlertsSection> {
    const capped = Math.min(Math.max(limit, 1), MAX_ALERTS);
    const { value, error } = await runSection("latestAlerts", async () => {
      const url = await root._urls.url("notification");
      const res = await ctx.dataSources.notification.getNotifications(
        url,
        ctx.headers.orgId,
        ctx.headers.userId
      );
      return res.notifications.slice(0, capped);
    });
    return { notifications: value, error };
  }

  @FieldResolver(() => KpisSection)
  async kpis(
    @Root() root: OverviewRoot,
    @Ctx() ctx: AppContext
  ): Promise<KpisSection> {
    // Phase 4: latest network KPIs polled from the metric service (closes
    // backend gap #5). Console polls this selection — no subscriptions in v1.
    const { value, error } = await runSection("kpis", async () => {
      // NB: the metric system's service-discovery name is "metrics"
      // (getSystemNameByService), not SUB_GRAPHS' "metric".
      const url = await root._urls.url("metrics");
      return fetchLatestKpis(ctx.dataSources.metric, url, NETWORK_KPI_KEYS);
    });
    return { metrics: value, error };
  }
}
