/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { SimPoolStatsDto } from "./types";

@Resolver()
export class GetSimPoolStatsResolver {
  @Query(() => SimPoolStatsDto)
  @UseMiddleware(Authentication)
  async getSimPoolStats(
    @Arg("type") type: string,
    @Ctx() ctx: Context
  ): Promise<SimPoolStatsDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.getSimPoolStats(type);
  }
}
