/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { GetReportsDto, GetReportsInputDto } from "./types";

@Resolver()
export class GetReportsResolver {
  @Query(() => GetReportsDto)
  async getReports(
    @Arg("data") data: GetReportsInputDto,
    @Ctx() ctx: AppContext
  ): Promise<GetReportsDto> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("billing");
    return await dataSources.billing.getReports(baseURL, data);
  }
}
