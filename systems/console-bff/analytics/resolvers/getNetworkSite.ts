/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { AnalyticsSiteInput } from "./types/business";
import { NetworkSiteDto } from "./types/network";

@Resolver()
export class GetNetworkSiteResolver {
  @Query(() => NetworkSiteDto)
  async getNetworkSite(
    @Arg("data") data: AnalyticsSiteInput,
    @Ctx() ctx: AppContext
  ): Promise<NetworkSiteDto> {
    const baseURL = await ctx.urls.url("analytics");
    return ctx.dataSources.analytics.getNetworkSite(baseURL, data);
  }
}
