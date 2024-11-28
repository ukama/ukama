/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { AddPackageSimResDto, AddPackagesToSimInputDto } from "./types";

@Resolver()
export class AddPackagesToSimResolver {
  @Mutation(() => [AddPackageSimResDto])
  async addPackagesToSim(
    @Arg("data") data: AddPackagesToSimInputDto,
    @Ctx() ctx: Context
  ): Promise<AddPackageSimResDto[]> {
    const { dataSources, baseURL } = ctx;
    return dataSources.dataSource.AddPackagesToSim(baseURL, data);
  }
}
