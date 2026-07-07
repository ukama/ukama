/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { SalesOverviewDto } from "./types/business";
import { AnalyticsWindowInput } from "./types/shared";

/** Parse a numeric-string money field; non-numeric → 0. */
const money = (s: string | undefined | null): number => {
  const n = Number(s);
  return Number.isFinite(n) ? n : 0;
};

@Resolver()
export class GetSalesOverviewResolver {
  @Query(() => SalesOverviewDto)
  async getSalesOverview(
    @Arg("data") data: AnalyticsWindowInput,
    @Ctx() ctx: AppContext
  ): Promise<SalesOverviewDto> {
    const baseURL = await ctx.urls.url("analytics");
    const overview = await ctx.dataSources.analytics.getSalesOverview(
      baseURL,
      data
    );

    // "Pending revenue" = real outstanding amount from the payment service
    // (analytics has no pending source). Per payment, outstanding = amount not
    // yet collected: a paid payment (paidAt set) contributes 0, an unpaid one
    // contributes amount − depositedAmount. Omitted on failure — never
    // fabricated. See systems/console-bff/BACKEND-GAPS.md.
    if (!overview.kpis.some(k => k.key === "revenue_pending")) {
      try {
        const paymentsURL = await ctx.urls.url("payments");
        const { payments } = await ctx.dataSources.payment.getPayments(
          paymentsURL,
          {}
        );
        const pending = payments.reduce(
          (sum, p) =>
            p.paidAt
              ? sum
              : sum + Math.max(0, money(p.amount) - money(p.depositedAmount)),
          0
        );
        const value = Math.round(pending * 100) / 100;
        overview.kpis = [
          ...overview.kpis,
          { key: "revenue_pending", value, formatted: String(value) },
        ];
      } catch {
        // Payment service unavailable — leave revenue_pending absent.
      }
    }

    return overview;
  }
}
