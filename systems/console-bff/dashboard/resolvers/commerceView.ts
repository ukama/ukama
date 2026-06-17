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

import { PaymentsDto } from "../../payment/resolver/types";
import type { AppContext } from "../../server/context";
import { ServiceUrlResolver } from "../baseUrls";
import { derivePlanStats, summarizeRevenue } from "../derive";
import { runSection } from "../section";
import {
  BalanceSection,
  CommerceView,
  InvoicesSection,
  PlanStatsSection,
  RevenueSection,
} from "./types";

const MAX_INVOICES = 100;

type CommerceRoot = CommerceView & {
  _urls: ServiceUrlResolver;
  /** payments memo — `revenue` and `plans` share one upstream list. */
  _payments?: Promise<PaymentsDto>;
};

/**
 * Business-lens commerce composite (plan §3.1). Serves: Business home
 * (revenue summary), Revenue screen, Packages screen (plans), Billing screen
 * (invoices + balance) — each through its own selection.
 */
@Resolver(() => CommerceView)
export class CommerceViewResolver {
  @Query(() => CommerceView)
  commerceView(
    @Ctx() ctx: AppContext,
    @Arg("networkId", { nullable: true }) networkId?: string
  ): CommerceView {
    return Object.assign(new CommerceView(), {
      networkId,
      _urls: new ServiceUrlResolver(ctx.headers.orgName),
    });
  }

  private fetchPayments(
    root: CommerceRoot,
    ctx: AppContext
  ): Promise<PaymentsDto> {
    if (!root._payments) {
      root._payments = root._urls
        .url("payments")
        .then(url => ctx.dataSources.payment.getPayments(url, {}));
    }
    return root._payments;
  }

  @FieldResolver(() => RevenueSection)
  async revenue(
    @Root() root: CommerceRoot,
    @Ctx() ctx: AppContext
  ): Promise<RevenueSection> {
    const { value, error } = await runSection("revenue", async () => {
      const payments = await this.fetchPayments(root, ctx);
      return summarizeRevenue(payments.payments);
    });
    return { ...value, error };
  }

  @FieldResolver(() => PlanStatsSection)
  async plans(
    @Root() root: CommerceRoot,
    @Ctx() ctx: AppContext
  ): Promise<PlanStatsSection> {
    const { value, error } = await runSection("plans", async () => {
      const [packagesUrl, payments] = await Promise.all([
        root._urls.url("package"),
        this.fetchPayments(root, ctx),
      ]);
      const packagesPromise = ctx.dataSources.package.getPackages(packagesUrl);
      // Attach counts need a network scope; without one they stay null.
      const simsPromise = root.networkId
        ? root._urls
            .url("subscriber")
            .then(url =>
              ctx.dataSources.subscriber.getSimsByNetwork(
                url,
                root.networkId as string
              )
            )
        : Promise.resolve(null);
      const [packages, sims] = await Promise.all([
        packagesPromise,
        simsPromise,
      ]);
      const simPackages = sims
        ? sims.sims.map(sim => ({
            packageId: sim.package?.package_id,
            isActive: sim.package?.is_active,
          }))
        : null;
      const plans = derivePlanStats(
        packages.packages,
        payments.payments,
        simPackages
      );
      const revenue = summarizeRevenue(payments.payments);
      const attachTotal = simPackages
        ? plans.reduce((acc, plan) => acc + (plan.attachCount ?? 0), 0)
        : null;
      return {
        plans,
        mrr: revenue.monthPaid,
        arpu:
          attachTotal && attachTotal > 0
            ? Math.round((revenue.monthPaid / attachTotal) * 100) / 100
            : null,
      };
    });
    return { ...value, error };
  }

  @FieldResolver(() => InvoicesSection)
  async invoices(
    @Root() root: CommerceRoot,
    @Ctx() ctx: AppContext,
    @Arg("limit", () => Int, { defaultValue: 20 }) limit: number
  ): Promise<InvoicesSection> {
    const capped = Math.min(Math.max(limit, 1), MAX_INVOICES);
    const { value, error } = await runSection("invoices", async () => {
      const url = await root._urls.url("billing");
      const res = await ctx.dataSources.billing.getReports(url, {
        networkId: root.networkId,
      });
      return res.reports.slice(0, capped);
    });
    return { reports: value, error };
  }

  @FieldResolver(() => BalanceSection)
  async balance(
    @Root() root: CommerceRoot,
    @Ctx() ctx: AppContext
  ): Promise<BalanceSection> {
    const { value, error } = await runSection("balance", async () => {
      const url = await root._urls.url("billing");
      const res = await ctx.dataSources.billing.getReports(url, {
        networkId: root.networkId,
        isPaid: false,
      });
      const unpaid = res.reports.filter(report => !report.isPaid);
      // Outstanding amount is derived from each invoice's raw report
      // (totalAmountCents, synced from the billing provider). Cents → major
      // units; non-numeric/missing totals contribute 0.
      const outstandingCents = unpaid.reduce((acc, report) => {
        const cents = Number(report.rawReport?.totalAmountCents);
        return acc + (Number.isFinite(cents) ? cents : 0);
      }, 0);
      return {
        outstandingCount: unpaid.length,
        latestUnpaidPeriod: unpaid[0]?.period ?? null,
        outstandingAmount: unpaid.length ? outstandingCents / 100 : null,
        currency: unpaid[0]?.rawReport?.currency ?? null,
      };
    });
    return { ...value, error };
  }
}
