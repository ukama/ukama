/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, FieldResolver, Query, Resolver, Root } from "type-graphql";

import type { AppContext } from "../../server/context";
import { ServiceUrlResolver } from "../baseUrls";
import { groupBy } from "../derive";
import { notImplementedSection, runSection } from "../section";
import {
  GapSection,
  PlanNameDto,
  PlansSection,
  SubscribersSection,
  SubscribersView,
} from "./types";

type SubscribersViewRoot = SubscribersView & { _urls: ServiceUrlResolver };

/**
 * Subscriber list composite (plan §3.1). Serves: Network customers,
 * Business customers, Customer-lens list — each selects different columns.
 */
@Resolver(() => SubscribersView)
export class SubscribersViewResolver {
  @Query(() => SubscribersView)
  subscribersView(
    @Arg("networkId") networkId: string,
    @Ctx() ctx: AppContext
  ): SubscribersView {
    return Object.assign(new SubscribersView(), {
      networkId,
      _urls: new ServiceUrlResolver(ctx.headers.orgName),
    });
  }

  @FieldResolver(() => SubscribersSection)
  async subscribers(
    @Root() root: SubscribersViewRoot,
    @Ctx() ctx: AppContext
  ): Promise<SubscribersSection> {
    const { value, error } = await runSection("subscribers", async () => {
      const url = await root._urls.url("subscriber");
      const [sims, subs] = await Promise.all([
        ctx.dataSources.subscriber.getSimsByNetwork(url, root.networkId),
        ctx.dataSources.subscriber.getSubscribersByNetwork(url, root.networkId),
      ]);
      const simsBySubscriber = groupBy(sims.sims, sim => sim.subscriberId);
      for (const sub of subs.subscribers) {
        sub.sim = simsBySubscriber.get(sub.uuid) ?? [];
      }
      return subs.subscribers;
    });
    return { subscribers: value, error };
  }

  @FieldResolver(() => PlansSection)
  async plans(
    @Root() root: SubscribersViewRoot,
    @Ctx() ctx: AppContext
  ): Promise<PlansSection> {
    const { value, error } = await runSection("plans", async () => {
      const url = await root._urls.url("package");
      const res = await ctx.dataSources.package.getPackages(url);
      return res.packages.map(
        (pkg): PlanNameDto => ({ packageId: pkg.uuid, name: pkg.name })
      );
    });
    return { plans: value, error };
  }

  @FieldResolver(() => GapSection)
  usage(): GapSection {
    // TODO(backend-gap): subscriber — batch usage endpoint
    // `/v1/usages?sim_ids=` + per-subscriber aggregation (today: one call per
    // sim) — unblocks: subscribersView.usage (docs/backend-gaps.md #2)
    return { error: notImplementedSection("usage").error };
  }
}
