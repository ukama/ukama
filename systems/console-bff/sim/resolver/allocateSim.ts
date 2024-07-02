/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { AllocateSimAPIDto, AllocateSimInputDto } from "./types";

@Resolver()
export class AllocateSimResolver {
  @Mutation(() => AllocateSimAPIDto)
  async allocateSim(
    @Arg("data") data: AllocateSimInputDto,
    @Ctx() ctx: Context
  ): Promise<AllocateSimAPIDto> {
    const { dataSources, baseURL } = ctx;
    return await dataSources.dataSource.allocateSim(baseURL, data);
  }
}
