/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { CBooleanResponse } from "../../common/types";
import type { AppContext } from "../../server/context";

@Resolver()
export class DeleteSiteResolver {
  @Mutation(() => CBooleanResponse)
  async deleteSite(
    @Arg("id") id: string,
    @Ctx() ctx: AppContext
  ): Promise<CBooleanResponse> {
    const baseURL = await ctx.urls.url("site");
    return ctx.dataSources.site.deleteSite(baseURL, id);
  }
}
