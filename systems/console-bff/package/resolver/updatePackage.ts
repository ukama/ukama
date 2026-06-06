/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { PackageDto, UpdatePackageInputDto } from "./types";

@Resolver()
export class UpdatePackageResolver {
  @Mutation(() => PackageDto)
  async updatePackage(
    @Arg("packageId") packageId: string,
    @Arg("data") data: UpdatePackageInputDto,
    @Ctx() ctx: AppContext
  ): Promise<PackageDto> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("package");
    return dataSources.package.updatePackage(baseURL, packageId, data);
  }
}
