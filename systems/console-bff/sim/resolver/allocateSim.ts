/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { AllocateSimAPIDto, AllocateSimInputDto } from "./types";

@Resolver()
export class AllocateSimResolver {
  @Mutation(() => AllocateSimAPIDto)
  async allocateSim(
    @Arg("data") data: AllocateSimInputDto,
    @Ctx() ctx: AppContext
  ): Promise<AllocateSimAPIDto> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("sim");

    let simToken: string | undefined;
    if (data.iccid) {
      simToken = await dataSources.sim.getTokenByIccid(baseURL, data.iccid);
    }

    return await dataSources.sim.allocateSim(baseURL, data, simToken);
  }
}
