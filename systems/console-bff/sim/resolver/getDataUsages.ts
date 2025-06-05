/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { logger } from "../../common/logger";
import { Context } from "../context";
import { SimDataUsages, SimDto, SimUsagesInputDto } from "./types";

@Resolver()
export class GetDataUsagesResolver {
  @Query(() => SimDataUsages)
  async getDataUsages(
    @Arg("data") data: SimUsagesInputDto,
    @Ctx() ctx: Context
  ): Promise<SimDataUsages> {
    try {
      const { dataSources, baseURL } = ctx;

      const [activeSims, inactiveSims] = await Promise.all([
        dataSources.dataSource.list(baseURL, {
          networkId: data.networkId,
          status: "active",
        }),
        dataSources.dataSource.list(baseURL, {
          networkId: data.networkId,
          status: "inactive",
        }),
      ]);

      const allSims = [
        ...(activeSims.sims || []),
        ...(inactiveSims.sims || []),
      ];
      const validSims = allSims.filter(
        (s: SimDto) =>
          s?.id && s?.iccid && s?.package?.id && s?.package?.startDate
      );

      logger.info(`Processing ${validSims.length} sims for usage data`);

      const usages = await Promise.all(
        validSims.map(async (sim: SimDto) => {
          try {
            return await dataSources.dataSource.getDataUsage(baseURL, {
              type: data.type,
              iccid: sim.iccid!,
              simId: sim.id!,
              from: sim.package!.startDate!,
              to: new Date().toISOString(),
            });
          } catch (error) {
            logger.warn(`Failed to get usage for SIM ${sim.id}: ${error}`);
            return { usage: "0", simId: sim.id! };
          }
        })
      );

      return { usages: usages.filter(usage => usage?.simId) };
    } catch (error) {
      logger.error(`Error fetching data usages: ${error}`);
      throw error;
    }
  }
}
