/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { AddPackageSimResDto, AddPackageToSimInputDto } from "./types";

@Resolver()
export class AddPackageToSimResolver {
  @Mutation(() => AddPackageSimResDto)
  @UseMiddleware(Authentication)
  async addPackageToSim(
    @Arg("data") data: AddPackageToSimInputDto,
    @Ctx() ctx: Context
  ): Promise<AddPackageSimResDto> {
    const { dataSources, baseURL } = ctx;
    return await dataSources.dataSource.addPackageToSim(baseURL, data);
  }
}
