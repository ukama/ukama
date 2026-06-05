/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { AddNetworkInputDto, NetworkDto } from "./types";

@Resolver()
export class AddNetworkResolver {
  @Mutation(() => NetworkDto)
  async addNetwork(
    @Arg("data") data: AddNetworkInputDto,
    @Ctx() ctx: AppContext
  ): Promise<NetworkDto> {
    const baseURL = await ctx.urls.url("network");
    return ctx.dataSources.network.addNetwork(baseURL, data);
  }
}
