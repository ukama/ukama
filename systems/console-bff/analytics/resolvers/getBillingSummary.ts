/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { BillingSummaryDto } from "./business.types";
import { AnalyticsWindowInput } from "./shared";

@Resolver()
export class GetBillingSummaryResolver {
  @Query(() => BillingSummaryDto)
  async getBillingSummary(
    @Arg("data") data: AnalyticsWindowInput,
    @Ctx() ctx: AppContext
  ): Promise<BillingSummaryDto> {
    const baseURL = await ctx.urls.url("analytics");
    return ctx.dataSources.analytics.getBillingSummary(baseURL, data);
  }
}
