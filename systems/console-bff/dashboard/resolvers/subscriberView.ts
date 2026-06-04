/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, FieldResolver, Query, Resolver, Root } from "type-graphql";

import type { AppContext } from "../../server/context";
import { SubscriberDto } from "../../subscriber/resolver/types";
import { ServiceUrlResolver } from "../baseUrls";
import { isPaidPayment } from "../derive";
import { notImplementedSection, runSection } from "../section";
import {
  GapSection,
  PlanNameDto,
  SubscriberBillingSection,
  SubscriberPlansSection,
  SubscriberSection,
  SubscriberView,
} from "./types";

const MAX_PAYMENTS = 20;

type SubscriberRoot = SubscriberView & {
  _urls: ServiceUrlResolver;
  /** subscriber core memo — `plans` and `billing` reuse it. */
  _subscriber?: Promise<SubscriberDto>;
};

/**
 * Customer-lens subscriber composite (plan §3.1). Serves: Customer home,
 * subscriber detail drawer (any lens).
 */
@Resolver(() => SubscriberView)
export class SubscriberViewResolver {
  @Query(() => SubscriberView)
  subscriberView(
    @Arg("subscriberId") subscriberId: string,
    @Ctx() ctx: AppContext
  ): SubscriberView {
    return Object.assign(new SubscriberView(), {
      subscriberId,
      _urls: new ServiceUrlResolver(ctx.headers.orgName),
    });
  }

  private fetchSubscriber(
    root: SubscriberRoot,
    ctx: AppContext
  ): Promise<SubscriberDto> {
    if (!root._subscriber) {
      root._subscriber = root._urls
        .url("subscriber")
        .then(url =>
          ctx.dataSources.subscriber.getSubscriber(url, root.subscriberId)
        );
    }
    return root._subscriber;
  }

  @FieldResolver(() => SubscriberSection)
  async subscriber(
    @Root() root: SubscriberRoot,
    @Ctx() ctx: AppContext
  ): Promise<SubscriberSection> {
    const { value, error } = await runSection("subscriber", () =>
      this.fetchSubscriber(root, ctx)
    );
    return { subscriber: value, error };
  }

  @FieldResolver(() => SubscriberPlansSection)
  async plans(
    @Root() root: SubscriberRoot,
    @Ctx() ctx: AppContext
  ): Promise<SubscriberPlansSection> {
    const { value, error } = await runSection("plans", async () => {
      const [subscriber, packagesUrl] = await Promise.all([
        this.fetchSubscriber(root, ctx),
        root._urls.url("package"),
      ]);
      const packages = await ctx.dataSources.package.getPackages(packagesUrl);
      const namesById = new Map(
        packages.packages.map(pkg => [pkg.uuid, pkg.name])
      );
      const activeIds = new Set(
        (subscriber.sim ?? [])
          .map(sim => sim.package)
          .filter(pkg => pkg?.is_active)
          .map(pkg => pkg?.package_id as string)
      );
      return Array.from(
        activeIds,
        (packageId): PlanNameDto => ({
          packageId,
          name: namesById.get(packageId) ?? packageId,
        })
      );
    });
    return { plans: value, error };
  }

  @FieldResolver(() => SubscriberBillingSection)
  async billing(
    @Root() root: SubscriberRoot,
    @Ctx() ctx: AppContext
  ): Promise<SubscriberBillingSection> {
    const { value, error } = await runSection("billing", async () => {
      const [subscriber, paymentsUrl] = await Promise.all([
        this.fetchSubscriber(root, ctx),
        root._urls.url("payments"),
      ]);
      const res = await ctx.dataSources.payment.getPayments(paymentsUrl, {});
      // TODO(backend-gap): payments — payer/subscriber filter on
      // /v1/payments (today: filter by payerEmail in the BFF) — unblocks:
      // subscriberView.billing at scale
      return res.payments
        .filter(
          payment =>
            payment.payerEmail &&
            payment.payerEmail.toLowerCase() === subscriber.email.toLowerCase()
        )
        .sort((a, b) => Number(isPaidPayment(b)) - Number(isPaidPayment(a)))
        .slice(0, MAX_PAYMENTS);
    });
    return { payments: value, error };
  }

  @FieldResolver(() => GapSection)
  usage(): GapSection {
    // TODO(backend-gap): subscriber — batch usage endpoint + per-subscriber
    // aggregation — unblocks: subscriberView.usage (docs/backend-gaps.md #2)
    return { error: notImplementedSection("usage").error };
  }
}
