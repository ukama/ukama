/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { GetKpiValuesDto, KpiValuesInput } from "./types/kpi";

@Resolver()
export class GetKpiValuesResolver {
  @Query(() => GetKpiValuesDto)
  async getKpiValues(
    @Arg("data") data: KpiValuesInput,
    @Ctx() ctx: AppContext
  ): Promise<GetKpiValuesDto> {
    const baseURL = await ctx.urls.url("analytics");
    return ctx.dataSources.analytics.getKpiValues(
      baseURL,
      data,
      ctx.headers.currency
    );
  }
}
