/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { CBooleanResponse } from "../../common/types";
import type { AppContext } from "../../server/context";
import { SetSiteInputDto } from "./types";

@Resolver()
export class SetSiteResolver {
  @Mutation(() => CBooleanResponse)
  async setSite(
    @Arg("data") data: SetSiteInputDto,
    @Ctx() ctx: AppContext
  ): Promise<CBooleanResponse> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("controller");
    return dataSources.controller.setSite(baseURL, data);
  }
}
