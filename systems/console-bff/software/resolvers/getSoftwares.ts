/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { GetSoftwaresInput, Softwares } from "./types";

@Resolver()
export class GetSoftwares {
  @Query(() => Softwares)
  async getSoftwares(
    @Ctx() ctx: AppContext,
    @Arg("data") data: GetSoftwaresInput
  ): Promise<Softwares> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("software");
    return await dataSources.software.getSoftwares(baseURL, data);
  }
}
