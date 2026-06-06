/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { logger } from "../../common/logger";
import { mapWithConcurrencySettled } from "../../common/utils/concurrency";
import type { AppContext } from "../../server/context";
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
    @Ctx() ctx: AppContext
  ): Promise<AddPackagesSimResDto> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("sim");

    // Parallel with bounded concurrency; report per-package outcome instead
    // of failing the whole mutation on the first error.
    const results = await mapWithConcurrencySettled(
      data.packages,
      packageInfo =>
        dataSources.sim.addPackageToSim(
          baseURL,
          data.sim_id,
          packageInfo.package_id,
          packageInfo.start_date
        )
    );

    const packages: AddPackagSimResDto[] = results.map((result, i) => {
      const packageId = data.packages[i].package_id;
      if (result.status === "rejected") {
        logger.error(
          `Error adding package to sim: ${packageId}: ${result.reason}`
        );
        return {
          packageId,
          success: false,
          error: "Failed to add package to sim",
        };
      }
      return { packageId, success: true };
    });

    return {
      packages,
    };
  }
}
