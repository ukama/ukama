/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { logger } from "../../common/logger";
import { Context } from "../context";
import {
  AddPackagSimResDto,
  AddPackagesSimResDto,
  AddPackagesToSimInputDto,
} from "./types";

@Resolver()
export class AddPackagesToSimResolver {
  @Mutation(() => AddPackagesSimResDto)
  async addPackagesToSim(
    @Arg("data") data: AddPackagesToSimInputDto,
    @Ctx() ctx: Context
  ): Promise<AddPackagesSimResDto> {
    const { dataSources, baseURL } = ctx;
    const pacakgesId: AddPackagSimResDto[] = [];
    for (const packageInfo of data.packages) {
      try {
        await dataSources.dataSource.addPackageToSim(
          baseURL,
          data.sim_id,
          packageInfo.package_id,
          packageInfo.start_date
        );
        pacakgesId.push({
          packageId: packageInfo.package_id,
        });
      } catch (error) {
        logger.error(`Error adding package to sim: ${packageInfo.package_id} `);
        throw new Error("Failed to add package to sim");
      }
    }
    return {
      packages: pacakgesId,
    };
  }
}
