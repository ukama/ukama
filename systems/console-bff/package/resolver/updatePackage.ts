/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { PackageDto, UpdatePackageInputDto } from "./types";

@Resolver()
export class UpdatePackageResolver {
  @Mutation(() => PackageDto)
  async updatePackage(
    @Arg("packageId") packageId: string,
    @Arg("data") data: UpdatePackageInputDto,
    @Ctx() ctx: Context
  ): Promise<PackageDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.updatePackage(packageId, data);
  }
}
