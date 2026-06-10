/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { RefreshInput, RefreshResultDto } from "./types/collector";

@Resolver()
export class RefreshAnalyticsResolver {
  @Mutation(() => RefreshResultDto)
  async refreshAnalytics(
    @Arg("data") data: RefreshInput,
    @Ctx() ctx: AppContext
  ): Promise<RefreshResultDto> {
    const baseURL = await ctx.urls.url("analytics");
    return ctx.dataSources.analytics.refreshAnalytics(baseURL, data);
  }
}
