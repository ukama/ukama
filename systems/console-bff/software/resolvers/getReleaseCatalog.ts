/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { GetReleaseCatalogInput, ReleaseCatalog } from "./types";

@Resolver()
export class GetReleaseCatalog {
  @Query(() => ReleaseCatalog)
  async getReleaseCatalog(
    @Ctx() ctx: AppContext,
    @Arg("data") data: GetReleaseCatalogInput
  ): Promise<ReleaseCatalog> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("software");
    return dataSources.software.getReleaseCatalog(baseURL, data);
  }
}
