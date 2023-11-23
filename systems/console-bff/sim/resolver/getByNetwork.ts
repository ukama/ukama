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
import { GetSimByNetworkInputDto, SimDetailsDto } from "./types";

@Resolver()
export class GetSimByNetworkResolver {
  @Query(() => SimDetailsDto)
  @UseMiddleware(Authentication)
  async getSim(
    @Arg("data") data: GetSimByNetworkInputDto,
    @Ctx() ctx: Context
  ): Promise<SimDetailsDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.getSimByNetworkId(data);
  }
}
