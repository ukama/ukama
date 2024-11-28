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
    const addedPackageIds: AddPackageSimResDto[] = [];

    for (const packageInfo of data.packages) {
      try {
        await dataSources.dataSource.AddPackagesToSim(baseURL, {
          sim_id: data.sim_id,
          packages: [packageInfo],
        });

        addedPackageIds.push({
          packageId: packageInfo.package_id,
        });
      } catch (error) {
        console.error(`Failed to add package ${packageInfo.package_id}`, error);
      }
    }

    return addedPackageIds;
  }
}
