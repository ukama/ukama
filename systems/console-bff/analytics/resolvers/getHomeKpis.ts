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
    return ctx.dataSources.analytics.getHomeKpis(baseURL, data);
  }
}
