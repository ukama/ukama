/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { Apps, GetAppsInputDto } from "./types";

@Resolver()
export class GetApps {
  @Query(() => Apps, { nullable: true })
  async getApps(
    @Ctx() ctx: AppContext,
    @Arg("data") data: GetAppsInputDto
  ): Promise<Apps> {
    // health resolves to the node gateway (nodeGwIp:nodeGwPort) via isForNodeGw.
    const baseURL = await ctx.urls.url("health");
    return ctx.dataSources.health.getApps(baseURL, data);
  }
}
