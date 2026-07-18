/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { PromoteReleaseInputDto, PromoteReleaseResponse } from "./types";

@Resolver()
export class PromoteReleaseResolver {
  @Mutation(() => PromoteReleaseResponse)
  async promoteRelease(
    @Arg("data") data: PromoteReleaseInputDto,
    @Ctx() ctx: AppContext
  ): Promise<PromoteReleaseResponse> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("software");
    return dataSources.software.promoteRelease(baseURL, data);
  }
}
