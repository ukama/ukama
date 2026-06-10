/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { GetResourceLockInputDto, ResourceLockDto } from "./types";

@Resolver()
export class GetResourceLockResolver {
  @Query(() => ResourceLockDto)
  async getResourceLock(
    @Arg("data") data: GetResourceLockInputDto,
    @Ctx() ctx: AppContext
  ): Promise<ResourceLockDto> {
    const baseURL = await ctx.urls.url("operation");
    return ctx.dataSources.operation.getResourceLock(baseURL, data.resourceKey);
  }
}
