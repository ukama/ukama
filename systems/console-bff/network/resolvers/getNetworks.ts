/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { NetworksResDto } from "./types";

@Resolver()
export class GetNetworksResolver {
  @Query(() => NetworksResDto)
  async getNetworks(@Ctx() ctx: AppContext): Promise<NetworksResDto> {
    const baseURL = await ctx.urls.url("network");
    return ctx.dataSources.network.getNetworks(baseURL);
  }
}
